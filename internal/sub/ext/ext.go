package ext

import (
	"context"

	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"github.com/etkecc/buscarron/internal/config"
	"github.com/etkecc/buscarron/internal/sub/ext/common"
	"github.com/etkecc/buscarron/internal/sub/ext/etkecc"
)

// Extension is form extension interface
type Extension interface {
	Execute(context.Context, common.Validator, *config.Form, map[string]string) (htmlResponse, matrixMessage string, files []*mautrix.ReqUploadMedia)
	Validate(context.Context, common.Validator, *config.Form, map[string]string) error
}

// EmailSender interface
type EmailSender interface {
	Send(context.Context, *postmark.Email) error
}

// New extensions map
func New(pm EmailSender) map[string]Extension {
	return map[string]Extension{
		"root":         NewRoot(),
		"confirmation": NewConfirmation(pm),
		"etkecc":       etkecc.New(pm),
	}
}
