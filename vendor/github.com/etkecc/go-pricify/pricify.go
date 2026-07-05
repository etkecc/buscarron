package pricify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/etkecc/go-kit/httpclient"
)

const dataURL = "https://etke.cc/order/components.json"

var (
	// go-kit single-host-tuned client, reused for the catalog fetch and the archive host.
	httpClient = httpclient.NewSingleHost()

	// userAgent resolves go-pricify's own imported version once from build info; v0.0.0 in dev.
	userAgent = sync.OnceValue(func() string {
		v := "v0.0.0"
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, dep := range info.Deps {
				if dep.Path == "github.com/etkecc/go-pricify" {
					v = dep.Version
					break
				}
			}
		}
		if v == "" || v == "(devel)" { // local replace / unstamped build
			v = "v0.0.0"
		}
		return "Go-Pricify-client/" + v
	})
)

// New price data, always returns cache (if available) on error
func New(ctx context.Context, uriOverride ...string) (*Data, error) {
	uri := dataURL
	if len(uriOverride) > 0 {
		uri = uriOverride[0]
	}
	source, err := load(ctx, uri)
	if err != nil {
		return getCache(), err
	}
	if source.ArchiveURL != "" {
		archiveSource, err := load(ctx, source.ArchiveURL)
		if err == nil {
			source.append(archiveSource)
		}
	}

	return convertToData(source), nil
}

func load(ctx context.Context, uri string) (*sourceModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// drain the dead body so httpclient pools the conn instead of torching it.
		_, _ = io.Copy(io.Discard, resp.Body) //nolint:errcheck // best-effort
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	sourceb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	source, err := parseSource(sourceb)
	if err != nil {
		return nil, err
	}
	return source, nil
}
