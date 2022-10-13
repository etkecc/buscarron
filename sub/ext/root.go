package ext

import (
	"sort"

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
func (e *root) Execute(_ common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	fields := e.sort(data)
	out := "**New " + form.Name + "**"
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
	out += "\n___\n"

	return out, []*mautrix.ReqUploadMedia{}
}

func (e *root) sort(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
