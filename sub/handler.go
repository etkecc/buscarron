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
	"gitlab.com/etke.cc/go/logger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext"
)

// Sender interface to send messages
type Sender interface {
	Send(id.RoomID, string)
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
	log         *logger.Logger
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
func NewHandler(forms map[string]*config.Form, vs map[string]Validator, pm EmailSender, sender Sender, loglevel string) *Handler {
	h := &Handler{
		redirectTpl: template.Must(template.New("redirect").Parse(redirect)),
		sanitizer:   bluemonday.StrictPolicy(),
		sender:      sender,
		forms:       forms,
		log:         logger.New("sub.", loglevel),
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
		h.log.Warn("submission attempt to the %s form (does not exist)", name)
		return "", ErrNotFound
	}
	v, ok := h.vs[name]
	if !ok {
		h.log.Warn("submission attempt to the %s form (validator does not exists)", name)
		return "", ErrNotFound
	}

	if err := r.ParseForm(); err != nil {
		h.log.Error("cannot parse a submission to the %s form: %v", name, err)
		return h.redirect(form.Redirect, nil), nil
	}

	data := make(map[string]string, len(r.PostForm))
	for key := range r.PostForm {
		data[key] = strings.TrimSpace(h.sanitizer.Sanitize(r.PostFormValue(key)))
	}

	if !v.Email(data["email"]) {
		h.log.Info("submission to the %s form marked as spam, reason: email", form.Name)
		return h.redirect(form.Redirect, data), ErrSpam
	}

	if !v.Domain(data["domain"]) {
		h.log.Info("submission to the %s form marked as spam, reason: domain", form.Name)
		return h.redirect(form.Redirect, data), ErrSpam
	}

	h.log.Info("submission attempt to the %s form by %v passed the tests", name, rID)

	text, files := h.generate(form, data)
	form.Lock()
	h.sender.Send(form.RoomID, text)
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
		h.log.Error("cannot parse redirect url template: %v", err)
	}
	if targetTpl != nil {
		err = targetTpl.Execute(&targetBytes, vars)
		if err != nil {
			h.log.Error("cannot execute redirect url template: %v", err)
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
		h.log.Error("cannot execute redirect template: %v", err)
	}

	return html.String()
}

// generate text and files
func (h *Handler) generate(form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	v := h.vs[form.Name]
	text, medias := h.ext["root"].Execute(v, form, data)

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

	return text, medias
}
