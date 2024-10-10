package ext

import (
	"context"
	"sort"
	"strings"

	"maunium.net/go/mautrix"

	"github.com/etkecc/buscarron/internal/config"
	"github.com/etkecc/buscarron/internal/sub/ext/common"
	"github.com/etkecc/buscarron/internal/utils"
)

type root struct{}

// NewRoot extension
func NewRoot() *root {
	return &root{}
}

// Execute extension
func (e *root) Execute(ctx context.Context, _ common.Validator, form *config.Form, data map[string]string) (htmlResponse, matrixMessage string, files []*mautrix.ReqUploadMedia) {
	span := utils.StartSpan(ctx, "sub.ext.root.Execute")
	defer span.Finish()

	defaultText := e.defaultText(form.Name, data)
	var out string
	if form.Text == "" {
		return "", defaultText, []*mautrix.ReqUploadMedia{}
	}

	out, err := common.ParseTemplate(form.Text, data)
	if err != nil {
		out = defaultText
	}
	out += "\n\n"

	files = []*mautrix.ReqUploadMedia{
		{
			Content:       strings.NewReader(defaultText),
			FileName:      "submission.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(defaultText)),
		},
	}

	return "", out, files
}

// Validate submission
func (e *root) Validate(_ context.Context, _ common.Validator, _ *config.Form, _ map[string]string) error {
	return nil
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
