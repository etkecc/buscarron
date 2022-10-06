package ext

import (
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
)

// Extension is form extension interface
type Extension interface {
	Execute(*config.Form, map[string]string) (string, []*mautrix.ReqUploadMedia)
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// Validator interface
type Validator interface {
	Domain(string, bool) bool
	DomainString(string) bool
	Email(string) bool
	A(string) bool
	CNAME(string) bool
	MX(string) bool
	GetBase(domain string) string
}

// New extensions map
func New(v Validator, pm EmailSender) map[string]Extension {
	return map[string]Extension{
		"root":         NewRoot(v),
		"confirmation": NewConfirmation(pm),
		"etkecc":       etkecc.New(v, pm),
	}
}
