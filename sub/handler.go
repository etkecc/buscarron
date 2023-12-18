package sub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/mattevans/postmark-go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/metrics"
	"gitlab.com/etke.cc/buscarron/sub/ext"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/linkpearl"
)

// Sender interface to send messages
type Sender interface {
	Send(id.RoomID, string, map[string]interface{}) id.EventID
	SendFile(id.RoomID, *mautrix.ReqUploadMedia, ...*event.RelatesTo)
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// Handler is an HTTP forms handler
type Handler struct {
	redirectTpl *template.Template
	sanitizer   *bluemonday.Policy
	sender      Sender
	mapping     map[string]func(r *http.Request) (map[string]string, error)
	forms       map[string]*config.Form
	log         *zerolog.Logger
	ext         map[string]ext.Extension
	vs          map[string]common.Validator
}

const redirect = "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='{{ .URL }}'\" /></head><body>Redirecting to <a href='{{ .URL }}'>{{ .URL }}</a>..."

var (
	// ErrNotFound when form not found
	ErrNotFound = errors.New("form not found")
	// ErrSpam returned when submission is spam
	ErrSpam = errors.New("spam submission")
)

// NewHandler creates new HTTP forms handler
func NewHandler(forms map[string]*config.Form, vs map[string]common.Validator, pm EmailSender, sender Sender, log *zerolog.Logger) *Handler {
	h := &Handler{
		redirectTpl: template.Must(template.New("redirect").Parse(redirect)),
		sanitizer:   bluemonday.StrictPolicy(),
		sender:      sender,
		forms:       forms,
		log:         log,
		ext:         ext.New(pm),
		vs:          vs,
	}
	h.initMapping()

	return h
}

func (h *Handler) initMapping() {
	h.mapping = map[string]func(r *http.Request) (map[string]string, error){
		"application/x-www-form-urlencoded": h.parseForm,
		"application/json":                  h.parseJSON,
	}
}

// GET request handler
func (h *Handler) GET(name string, _ *http.Request) (string, error) {
	form := h.forms[name]
	if form != nil {
		return h.redirect(form.Redirect, nil), nil
	}

	return "", ErrNotFound
}

// parseForm parses HTTP form (application/x-www-form-urlencoded)
func (h *Handler) parseForm(r *http.Request) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	data := make(map[string]string, len(r.PostForm))
	for key := range r.PostForm {
		data[key] = strings.TrimSpace(h.sanitizer.Sanitize(r.PostFormValue(key)))
	}

	return data, nil
}

func (h *Handler) parseJSON(r *http.Request) (map[string]string, error) {
	defer r.Body.Close()

	var rawData map[string]any
	err := json.NewDecoder(r.Body).Decode(&rawData)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string, len(rawData))
	for key, value := range rawData {
		data[key] = strings.TrimSpace(h.sanitizer.Sanitize(fmt.Sprintf("%+v", value)))
	}
	return data, nil
}

// POST request handler
func (h *Handler) POST(name string, r *http.Request) (string, error) {
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

	ctype := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))
	parser, ok := h.mapping[ctype]
	if !ok {
		h.log.Warn().Str("name", name).Str("content-type", ctype).Msg("form parser not found")
		return "", ErrNotFound
	}

	data, err := parser(r)
	if err != nil {
		return h.redirect(form.RejectRedirect, data), err
	}

	if !v.Email(data["email"]) {
		h.log.Info().Str("name", form.Name).Str("reason", "email").Msg("submission to the form marked as spam")
		return h.redirect(form.RejectRedirect, data), ErrSpam
	}

	if !v.Domain(data["domain"]) {
		h.log.Info().Str("name", form.Name).Str("reason", "domain").Msg("submission to the form marked as spam")
		return h.redirect(form.RejectRedirect, data), ErrSpam
	}

	metrics.Submission(form.Name)
	h.log.Info().Str("name", form.Name).Msg("submission to the form passed the tests")

	text, files := h.generate(form, data)
	attrs := map[string]interface{}{}
	if data["email"] != "" {
		attrs["email"] = data["email"]
	}
	if data["domain"] != "" {
		attrs["domain"] = data["domain"]
	}

	form.Lock()
	eventID := h.sender.Send(form.RoomID, text, attrs)
	var relates *event.RelatesTo
	if eventID != "" {
		relates = linkpearl.RelatesTo(eventID)
	}
	for _, file := range files {
		h.sender.SendFile(form.RoomID, file, relates)
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
