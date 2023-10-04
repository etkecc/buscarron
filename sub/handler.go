package sub

import (
	"bytes"
	"errors"
	"html/template"
	"net"
	"net/http"
	"strings"

	"github.com/mattevans/postmark-go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/metrics"
	"gitlab.com/etke.cc/buscarron/sub/ext"
)

// Sender interface to send messages
type Sender interface {
	Send(id.RoomID, string, map[string]interface{})
	SendFile(id.RoomID, *mautrix.ReqUploadMedia)
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// Validator interface
type Validator interface {
	Domain(string) bool
	DomainString(string) bool
	Email(string, ...net.IP) bool
	A(string) bool
	CNAME(string) bool
	MX(string) bool
	GetBase(domain string) string
}

// Handler is an HTTP forms handler
type Handler struct {
	redirectTpl *template.Template
	sanitizer   *bluemonday.Policy
	sender      Sender
	forms       map[string]*config.Form
	log         *zerolog.Logger
	ext         map[string]ext.Extension
	vs          map[string]Validator
}

const redirect = "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='{{ .URL }}'\" /></head><body>Redirecting to <a href='{{ .URL }}'>{{ .URL }}</a>..."

var (
	// ErrNotFound when form not found
	ErrNotFound = errors.New("form not found")
	// ErrSpam returned when submission is spam
	ErrSpam = errors.New("spam submission")
)

// NewHandler creates new HTTP forms handler
func NewHandler(forms map[string]*config.Form, vs map[string]Validator, pm EmailSender, sender Sender, log *zerolog.Logger) *Handler {
	h := &Handler{
		redirectTpl: template.Must(template.New("redirect").Parse(redirect)),
		sanitizer:   bluemonday.StrictPolicy(),
		sender:      sender,
		forms:       forms,
		log:         log,
		ext:         ext.New(pm),
		vs:          vs,
	}

	return h
}

// GET request handler
func (h *Handler) GET(name string, _ *http.Request) (string, error) {
	form := h.forms[name]
	if form != nil {
		return h.redirect(form.Redirect, nil), nil
	}

	return "", ErrNotFound
}

// POST request handler
func (h *Handler) POST(rID, name string, r *http.Request) (string, error) {
	form, ok := h.forms[name]
	if !ok {
		h.log.Warn().Str("name", name).Msg("submission attempt to a nonexistent form")
		return "", ErrNotFound
	}
	v, ok := h.vs[name]
	if !ok {
		h.log.Warn().Str("name", name).Msg("submission attempt to the a nonexistent form (validator does not exists)")
		return "", ErrNotFound
	}

	if err := r.ParseForm(); err != nil {
		h.log.Error().Str("name", name).Err(err).Msg("cannot parse a submission to the form")
		return h.redirect(form.Redirect, nil), nil
	}

	data := make(map[string]string, len(r.PostForm))
	for key := range r.PostForm {
		data[key] = strings.TrimSpace(h.sanitizer.Sanitize(r.PostFormValue(key)))
	}

	if !v.Email(data["email"]) {
		h.log.Info().Str("name", form.Name).Str("reason", "email").Msg("submission to the form marked as spam")
		return h.redirect(form.Redirect, data), ErrSpam
	}

	if !v.Domain(data["domain"]) {
		h.log.Info().Str("name", form.Name).Str("reason", "domain").Msg("submission to the form marked as spam")
		return h.redirect(form.Redirect, data), ErrSpam
	}

	metrics.Submission(form.Name)
	h.log.Info().Str("name", form.Name).Str("id", rID).Msg("submission to the form passed the tests")

	text, files := h.generate(form, data)
	attrs := map[string]interface{}{}
	if data["email"] != "" {
		attrs["email"] = data["email"]
	}
	if data["domain"] != "" {
		attrs["domain"] = data["domain"]
	}

	form.Lock()
	h.sender.Send(form.RoomID, text, attrs)
	for _, file := range files {
		h.sender.SendFile(form.RoomID, file)
	}
	form.Unlock()

	return h.redirect(form.Redirect, data), nil
}

func (h *Handler) redirect(target string, vars map[string]string) string {
	var html bytes.Buffer
	var targetBytes bytes.Buffer
	targetTpl, err := template.New("target").Parse(target)
	if err != nil {
		h.log.Error().Err(err).Msg("cannot parse redirect url template")
	}
	if targetTpl != nil {
		err = targetTpl.Execute(&targetBytes, vars)
		if err != nil {
			h.log.Error().Err(err).Msg("cannot execute redirect url template")
		} else {
			target = targetBytes.String()
		}
	}

	data := struct {
		URL string
	}{
		URL: target,
	}
	err = h.redirectTpl.Execute(&html, data)
	if err != nil {
		h.log.Error().Err(err).Msg("cannot execute redirect template")
	}

	return html.String()
}

// generate text and files
func (h *Handler) generate(form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	v := h.vs[form.Name]
	medias := []*mautrix.ReqUploadMedia{}
	text, rmedias := h.ext["root"].Execute(v, form, data)

	for _, extension := range form.Extensions {
		if extension == "" {
			continue
		}
		e, ok := h.ext[extension]
		if !ok || e == nil {
			continue
		}

		etext, emedias := e.Execute(v, form, data)
		text += etext
		medias = append(medias, emedias...)
	}
	medias = append(medias, rmedias...) // add submission.md at the end

	return text, medias
}
