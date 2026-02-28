package psd

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

const (
	// DefaultUserAgent is the default user agent used by the client
	DefaultUserAgent = "Go-PSD-client/v0.0.0"
	// DefaultMaxRetries is the default number of retries for failed requests
	DefaultMaxRetries = 3
	// DefaultMaxIdleConns is the maximum number of idle connections across all hosts.
	DefaultMaxIdleConns = 256
	// DefaultMaxIdleConnsPerHost is the maximum idle connections per host.
	DefaultMaxIdleConnsPerHost = 256
	// DefaultMaxConnsPerHost is the maximum total connections per host.
	DefaultMaxConnsPerHost = 256
	// DefaultIdleConnTimeout controls how long idle keep-alive connections stay around.
	DefaultIdleConnTimeout = 90 * time.Second
)

// Client is the client to interact with the Prometheus Service Discovery HTTP API
type Client struct {
	url        *url.URL
	baseURL    string
	baseNode   string
	baseHeader http.Header
	basicAuth  string
	cache      *expirable.LRU[string, cacheValue]
	client     *http.Client
	login      string
	password   string
	userAgent  string
	maxRetries int
}

type cacheValue struct {
	data     []byte
	etag     string
	cachedAt string
}

// NewClient returns a new PSD client
func NewClient(baseURL, login, password string, userAgentOverride ...string) *Client {
	uri, err := url.Parse(baseURL)
	if err != nil || login == "" || password == "" {
		return &Client{client: http.DefaultClient}
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &Client{client: http.DefaultClient}
	}
	transport := defaultTransport.Clone()
	transport.MaxIdleConns = DefaultMaxIdleConns
	transport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	transport.MaxConnsPerHost = DefaultMaxConnsPerHost
	transport.IdleConnTimeout = DefaultIdleConnTimeout
	transport.ForceAttemptHTTP2 = true
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ResponseHeaderTimeout = 10 * time.Second
	transport.ExpectContinueTimeout = 1 * time.Second

	cache := expirable.NewLRU[string, cacheValue](1000, nil, 5*time.Minute)
	userAgent := DefaultUserAgent
	if len(userAgentOverride) > 0 {
		userAgent = userAgentOverride[0]
	}

	baseHeader := make(http.Header, 4)
	baseHeader.Set("Content-Type", "application/json")
	baseHeader.Set("User-Agent", userAgent)
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
	baseHeader.Set("Authorization", basicAuth)

	baseURLStr := strings.TrimRight(uri.String(), "/")
	return &Client{
		url:        uri,
		baseURL:    baseURLStr,
		baseNode:   baseURLStr + "/node/",
		baseHeader: baseHeader,
		basicAuth:  basicAuth,
		cache:      cache,
		client:     &http.Client{Transport: transport},
		login:      login,
		password:   password,
		userAgent:  userAgent,
		maxRetries: DefaultMaxRetries,
	}
}

// buildURL constructs the URL for the given identifier and optional job override
func (p *Client) buildURL(identifier string, jobOverride ...string) string {
	if len(jobOverride) == 0 {
		return p.baseNode + identifier
	}
	job := jobOverride[0]

	var b strings.Builder
	b.Grow(len(p.baseURL) + 1 + len(job) + 1 + len(identifier))
	b.WriteString(p.baseURL)
	b.WriteByte('/')
	b.WriteString(job)
	b.WriteByte('/')
	b.WriteString(identifier)
	return b.String()
}

// do executes an HTTP request with per-attempt timeout and retry handling.
// It returns a cancel func that must be called after the response body is fully consumed.
func (p *Client) do(ctx context.Context, uri string, cached *cacheValue, currentAttempt ...int) (*http.Response, func(), error) {
	attempt := 0
	if len(currentAttempt) > 0 {
		attempt = currentAttempt[0]
	}

	headers := p.baseHeader.Clone()
	if cached != nil {
		headers.Set("If-Modified-Since", cached.cachedAt)
		headers.Set("If-None-Match", cached.etag)
	}
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	if _, ok := ctx.Deadline(); ok {
		cancel()
		childCtx = ctx
		cancel = func() {}
	}

	req, err := http.NewRequestWithContext(childCtx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	req.Header = headers

	resp, err := p.client.Do(req) //nolint:gosec // The URL is built from a trusted source
	if err != nil {
		cancel()
		return nil, nil, err
	}

	if resp.StatusCode == http.StatusGone { // not found, to distinguish from reverse proxy 404 error
		defer resp.Body.Close()
		cancel()
		return nil, nil, nil
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		defer resp.Body.Close()
		if attempt < p.maxRetries {
			wait := 1 * time.Second * time.Duration(attempt+1)
			cancel()
			select {
			case <-time.After(wait): // Exponential backoff
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			}
			return p.do(ctx, uri, cached, attempt+1)
		}

		err = fmt.Errorf("HTTP error: %s", resp.Status)
		cancel()
		return nil, nil, err
	}

	return resp, cancel, nil
}

// GetWithContext returns the list of targets for the given identifier using the given context
func (p *Client) GetWithContext(ctx context.Context, identifier string, jobOverride ...string) (Items, error) {
	if p.url == nil {
		return nil, nil
	}
	urlTarget := p.buildURL(identifier, jobOverride...)
	cachedData, cached := p.cache.Get(urlTarget)
	var cachedPtr *cacheValue
	if cached {
		cachedPtr = &cachedData
	}

	resp, cancel, err := p.do(ctx, urlTarget, cachedPtr)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	defer cancel()
	datab, err := p.responseData(resp, cached, cachedData, urlTarget)
	if err != nil {
		return nil, err
	}
	var psd []*Item
	err = json.Unmarshal(datab, &psd)
	if err != nil {
		return nil, err
	}

	return psd, nil
}

func (p *Client) responseData(resp *http.Response, cached bool, cachedData cacheValue, urlTarget string) ([]byte, error) {
	if resp.StatusCode == http.StatusNotModified && cached {
		return cachedData.data, nil
	}

	datab, err := readBody(resp.Body, resp.ContentLength)
	if err != nil {
		return nil, err
	}

	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		lastModified = time.Now().Format(http.TimeFormat)
	}
	p.cache.Add(urlTarget, cacheValue{
		data:     datab,
		etag:     resp.Header.Get("ETag"),
		cachedAt: lastModified,
	})

	return datab, nil
}

func readBody(r io.Reader, contentLength int64) ([]byte, error) {
	if contentLength > 0 && contentLength <= int64(int(^uint(0)>>1)) {
		var buf bytes.Buffer
		buf.Grow(int(contentLength))
		_, err := buf.ReadFrom(r)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return io.ReadAll(r)
}

// Get returns the list of targets for the given identifier
func (p *Client) Get(identifier string) (Items, error) {
	return p.GetWithContext(context.Background(), identifier)
}

// GetRaw returns the raw data for the given identifier, without using the cache or parsing the response.
// Do not use this method unless you need the raw data, as it does not parse the response
func (p *Client) GetRaw(ctx context.Context, identifier string, jobOverride ...string) ([]byte, error) {
	if p.url == nil {
		return nil, nil
	}
	urlTarget := p.buildURL(identifier, jobOverride...)
	resp, cancel, err := p.do(ctx, urlTarget, nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	defer cancel()

	return readBody(resp.Body, resp.ContentLength)
}
