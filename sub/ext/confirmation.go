package ext

import (
	"bytes"
	"text/template"

	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/buscarron/config"
	"maunium.net/go/mautrix"
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
func (e *confirmation) Execute(form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
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
	subject, err := e.parse(form.Confirmation.Subject, data)
	if err != nil {
		return "", []*mautrix.ReqUploadMedia{}
	}
	body, err := e.parse(form.Confirmation.Body, data)
	if err != nil {
		return "", []*mautrix.ReqUploadMedia{}
	}
	req := &postmark.Email{
		To:       email,
		Tag:      "confirmation",
		Subject:  subject,
		TextBody: body,
	}
	e.s.Send(req) // nolint // not ready to handle errors

	return "", []*mautrix.ReqUploadMedia{}
}

// parse template string
func (e *confirmation) parse(tplString string, data map[string]string) (string, error) {
	var result bytes.Buffer
	tpl, err := template.New("email").Parse(tplString)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
