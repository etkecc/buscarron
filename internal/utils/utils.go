package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

var sanitizer = bluemonday.StrictPolicy()

func Sanitize(input string) string {
	out := sanitizer.Sanitize(input)
	return strings.ReplaceAll(html.UnescapeString(out), "\"", "")
}
