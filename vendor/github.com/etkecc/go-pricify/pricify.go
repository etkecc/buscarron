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

	resp, err := http.Get(uri)
	if err != nil {
		return getCache(), err
	}

	sourceb, err := io.ReadAll(resp.Body)
	if err != nil {
		return getCache(), err
	}
	defer resp.Body.Close()

	source, err := parseSource(sourceb)
	if err != nil {
		return getCache(), err
	}
	return convertToData(source), nil
}
