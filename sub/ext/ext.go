package ext

import (
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
)

// Extension is form extension interface
type Extension interface {
	Execute(common.Validator, *config.Form, map[string]string) (string, []*mautrix.ReqUploadMedia)
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// New extensions map
func New(pm EmailSender) map[string]Extension {
	return map[string]Extension{
		"root":         NewRoot(),
		"confirmation": NewConfirmation(pm),
		"etkecc":       etkecc.New(pm),
	}
}
