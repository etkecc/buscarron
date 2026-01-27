package pricify

import (
	"io"
	"net/http"
)

const dataURL = "https://etke.cc/order/components.json"

// New price data, always returns cache (if available) on error
func New(uriOverride ...string) (*Data, error) {
	uri := dataURL
	if len(uriOverride) > 0 {
		uri = uriOverride[0]
	}
	source, err := load(uri)
	if err != nil {
		return getCache(), err
	}
	if source.ArchiveURL != "" {
		archiveSource, err := load(source.ArchiveURL)
		if err == nil {
			source.append(archiveSource)
		}
	}

	return convertToData(source), nil
}

func load(uri string) (*sourceModel, error) {
	resp, err := http.Get(uri) //nolint:gosec // intended
	if err != nil {
		return nil, err
	}

	sourceb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	source, err := parseSource(sourceb)
	if err != nil {
		return nil, err
	}
	return source, nil
}
