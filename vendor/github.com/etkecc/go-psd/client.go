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

	"github.com/etkecc/go-kit/httpclient"
	"github.com/goccy/go-json"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

// DefaultUserAgent is the default user agent used by the client
const DefaultUserAgent = "Go-PSD-client/v0.0.0"

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
		// 1s step, not the 200ms fleet default: a reloading SD backend gets room to breathe.
		client:    httpclient.NewSingleHost(httpclient.WithRetryDelayStep(1 * time.Second)),
		login:     login,
		password:  password,
		userAgent: userAgent,
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

// do issues one GET; httpclient owns retry and per-attempt timeout. A 200/304 is returned
// live for the caller to read and Close; a 410 (not-found) and any other non-2xx are
// terminal here, their bodies drained so the connection returns to the pool.
func (p *Client) do(ctx context.Context, uri string, cached *cacheValue) (*http.Response, error) {
	headers := p.baseHeader.Clone()
	if cached != nil {
		headers.Set("If-Modified-Since", cached.cachedAt)
		headers.Set("If-None-Match", cached.etag)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusGone { // not found, distinct from a reverse-proxy 404
		p.drain(resp)
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		p.drain(resp)
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return resp, nil
}

// drain empties and closes a body psd won't hand back, so httpclient's cancelOnClose fires
// and the connection returns to the pool instead of being torn down under a 5xx storm.
func (p *Client) drain(resp *http.Response) {
	_, _ = io.Copy(io.Discard, resp.Body) //nolint:errcheck // best-effort drain; a failure just costs the connection
	_ = resp.Body.Close()
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

	resp, err := p.do(ctx, urlTarget, cachedPtr)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
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
	resp, err := p.do(ctx, urlTarget, nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	return readBody(resp.Body, resp.ContentLength)
}
