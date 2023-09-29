package pricify

import (
	"io"
	"net/http"
)

// New price data
func New(uri string) (*Data, error) {
	resp, err := http.Get(uri)
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
	return convertToData(source), nil
}
