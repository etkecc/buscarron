package ext

import (
	"context"
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

type confirmation struct {
	s EmailSender
}

// NewConfirmation extension
func NewConfirmation(sender EmailSender) *confirmation {
	return &confirmation{s: sender}
}

// Execute extension
// nolint:unparam // interface constraints
func (e *confirmation) Execute(ctx context.Context, _ common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	span := sentry.StartSpan(ctx, "function", sentry.WithDescription("sub.ext.confirmation.Execute"))
	defer span.Finish()

	if e.s == nil {
		return "", []*mautrix.ReqUploadMedia{}
	}
	if form.Confirmation.Subject == "" && form.Confirmation.Body == "" {
		return "", []*mautrix.ReqUploadMedia{}
	}

	email, ok := data["email"]
	if !ok {
		return "", []*mautrix.ReqUploadMedia{}
	}
	subject, err := common.ParseTemplate(form.Confirmation.Subject, data)
	if err != nil {
		return "", []*mautrix.ReqUploadMedia{}
	}
	body, err := common.ParseTemplate(form.Confirmation.Body, data)
	if err != nil {
		return "", []*mautrix.ReqUploadMedia{}
	}
	req := &postmark.Email{
		To:       email,
		Tag:      "confirmation",
		Subject:  subject,
		TextBody: body,
	}

	// special case with nil interface
	if e.s == nil || (reflect.ValueOf(e.s).Kind() == reflect.Ptr && reflect.ValueOf(e.s).IsNil()) {
		return "", []*mautrix.ReqUploadMedia{}
	}

	e.s.Send(span.Context(), req) // nolint // not ready to handle errors
	return "", []*mautrix.ReqUploadMedia{}
}
