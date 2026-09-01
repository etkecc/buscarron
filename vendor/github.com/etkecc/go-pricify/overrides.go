package pricify

import (
	"strconv"
	"strings"
	"unicode"
)

// parseServerPriceOverrides parses comma-separated "size-region=price" entries into a lookup map; bad entries skipped.
func parseServerPriceOverrides(raw string) map[string]int {
	out := map[string]int{}
	for entry := range strings.SplitSeq(raw, ",") {
		key, value, found := strings.Cut(removeSpaces(entry), "=")
		if !found || key == "" {
			continue
		}
		if price, ok := parsePriceInt(value); ok {
			out[key] = price
		}
	}
	return out
}

// removeSpaces strips every whitespace rune so "size - region = price" collapses to "size-region=price".
func removeSpaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// parsePriceInt parses a major-unit price, tolerating float-style "15.00" input.
func parsePriceInt(raw string) (int, bool) {
	if price, err := strconv.Atoi(raw); err == nil {
		return price, true
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return int(f), true
}
