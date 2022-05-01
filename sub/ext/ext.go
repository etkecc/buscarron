package ext

import (
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
	"gitlab.com/etke.cc/buscarron/validator"
)

// Extension is form extension interface
type Extension interface {
	Execute(*config.Form, map[string]string) (string, []*mautrix.ReqUploadMedia)
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// New extensions map
func New(v *validator.V, pm EmailSender) map[string]Extension {
	return map[string]Extension{
		"root":         NewRoot(v),
		"confirmation": NewConfirmation(pm),
		"etkecc":       etkecc.New(v, pm),
	}
}
