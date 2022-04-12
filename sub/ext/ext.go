package ext

import (
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
	"gitlab.com/etke.cc/buscarron/validator"
)

// Extension is form extension interface
type Extension interface {
	Execute(string, map[string]string) (string, []*mautrix.ReqUploadMedia)
}

// New extensions map
func New(v *validator.V) map[string]Extension {
	return map[string]Extension{
		"root":   NewRoot(v),
		"etkecc": etkecc.New(v),
	}
}
