package ext

import (
	"context"
	"sort"
	"strings"

	"github.com/getsentry/sentry-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

type root struct{}

// NewRoot extension
func NewRoot() *root {
	return &root{}
}

// Execute extension
func (e *root) Execute(ctx context.Context, _ common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	span := sentry.StartSpan(ctx, "function", sentry.WithDescription("sub.ext.root.Execute"))
	defer span.Finish()

	defaultText := e.defaultText(form.Name, data)
	var out string
	if form.Text == "" {
		return defaultText, []*mautrix.ReqUploadMedia{}
	}

	out, err := common.ParseTemplate(form.Text, data)
	if err != nil {
		out = defaultText
	}
	out += "\n\n"

	files := []*mautrix.ReqUploadMedia{
		{
			Content:       strings.NewReader(defaultText),
			FileName:      "submission.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(defaultText)),
		},
	}

	return out, files
}

func (e *root) defaultText(name string, data map[string]string) string {
	fields := e.sort(data)
	out := "**New " + name + "**"
	if data["email"] != "" {
		out += " by " + data["email"] + "\n\n"
	} else {
		out += "\n\n"
	}

	for _, field := range fields {
		value := data[field]
		if value == "on" {
			value = "✅"
		}

		if value != "" {
			out += "* " + field + ": " + value + "\n"
		}
	}

	return out
}

func (e *root) sort(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
