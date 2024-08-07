package sub

import (
	"bytes"
	"context"
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
	"gitlab.com/etke.cc/buscarron/utils"
	"gitlab.com/etke.cc/go/redmine"
	"gitlab.com/etke.cc/linkpearl"
)

// Sender interface to send messages
type Sender interface {
	Send(context.Context, id.RoomID, string, map[string]any) id.EventID
	SendFile(context.Context, id.RoomID, *mautrix.ReqUploadMedia, ...*event.RelatesTo)
}

// EmailSender interface
type EmailSender interface {
	Send(context.Context, *postmark.Email) error
}

// Handler is an HTTP forms handler
type Handler struct {
	redirectTpl *template.Template
	sanitizer   *bluemonday.Policy
	sender      Sender
	mapping     map[string]func(r *http.Request) (map[string]string, error)
	forms       map[string]*config.Form
	rdm         *redmine.Redmine
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
func NewHandler(forms map[string]*config.Form, vs map[string]common.Validator, pm EmailSender, sender Sender, rdm *redmine.Redmine) *Handler {
	h := &Handler{
		redirectTpl: template.Must(template.New("redirect").Parse(redirect)),
		sanitizer:   bluemonday.StrictPolicy(),
		sender:      sender,
		forms:       forms,
		rdm:         rdm,
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

// SetSender sets sender
func (h *Handler) SetSender(sender Sender) {
	h.sender = sender
}

// GET request handler
func (h *Handler) GET(ctx context.Context, name string, _ *http.Request) (string, error) {
	span := utils.StartSpan(ctx, "sub.GET")
	defer span.Finish()

	form := h.forms[name]
	if form != nil {
		return h.redirect(span.Context(), form.Redirect, nil), nil
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
func (h *Handler) POST(ctx context.Context, name string, r *http.Request) (string, error) {
	log := zerolog.Ctx(ctx).With().Str("form", name).Logger()
	span := utils.StartSpan(ctx, "sub.POST")
	defer span.Finish()

	form, ok := h.forms[name]
	if !ok {
		log.Warn().Msg("submission attempt to a nonexistent form")
		return "", ErrNotFound
	}
	v, ok := h.vs[name]
	if !ok {
		log.Warn().Msg("submission attempt to an unknown form (validator does not exists)")
		return "", ErrNotFound
	}

	ctype := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))
	parser, ok := h.mapping[ctype]
	if !ok {
		log.Warn().Str("content-type", ctype).Msg("form parser not found")
		return "", ErrNotFound
	}

	data, err := parser(r)
	if err != nil {
		return h.redirect(span.Context(), form.RejectRedirect, data), err
	}

	if !v.Email(data["email"], "") {
		log.Info().Str("reason", "email").Msg("submission to the form marked as spam")
		return h.redirect(span.Context(), form.RejectRedirect, data), ErrSpam
	}

	if !v.Domain(data["domain"]) {
		log.Info().Str("reason", "domain").Msg("submission to the form marked as spam")
		return h.redirect(span.Context(), form.RejectRedirect, data), ErrSpam
	}

	if err := h.extValidate(span.Context(), form, data); err != nil {
		log.Info().Str("reason", "extension: "+err.Error()).Msg("submission to the form marked as spam")
		return h.redirect(span.Context(), form.RejectRedirect, data), ErrSpam
	}

	metrics.Submission(form.Name)
	log.Info().Msg("submission to the form passed the tests")

	delete(data, "issue_id") // remove issue_id from the data, to ensure it's not sent from the form, but generated by the system
	issueSubject := h.getIssueSubject(form.Name, data)
	issueID, err := h.rdm.NewIssue(issueSubject, "form", data["email"], "from buscarron")
	if err != nil {
		log.Error().Err(err).Msg("cannot create redmine issue")
	}

	attrs := map[string]any{}
	if issueID != 0 {
		data["issue_id"] = fmt.Sprintf("%d", issueID)
		attrs["issue_id"] = data["issue_id"]
	}

	log.Info().Msg("generating submission text and files")
	text, files := h.generate(span.Context(), form, data)
	log.Info().Msg("submission text and files have been generated")
	if data["email"] != "" {
		attrs["email"] = data["email"]
	}
	if data["domain"] != "" {
		attrs["base_domain"] = v.GetBase(data["domain"])
		attrs["domain"] = data["domain"]
	}

	go h.updateIssue(span.Context(), issueID, form.Name, text, files)

	log.Info().Msg("sending submission to the room")
	eventID := h.sender.Send(span.Context(), form.RoomID, text, attrs)
	var relates *event.RelatesTo
	if eventID != "" {
		relates = linkpearl.RelatesTo(eventID)
	}
	log.Info().Msg("submission has been sent to the room; sending files")
	for _, file := range files {
		h.sender.SendFile(span.Context(), form.RoomID, file, relates)
	}
	log.Info().Msg("files have been sent")

	return h.redirect(span.Context(), form.Redirect, data), nil
}

func (h *Handler) getIssueSubject(name string, data map[string]string) string {
	if data["domain"] != "" {
		return data["domain"]
	}
	if data["email"] != "" {
		return fmt.Sprintf("New %s submission by %s", name, data["email"])
	}
	return fmt.Sprintf("New %s submission", name)
}

func (h *Handler) updateIssue(ctx context.Context, issueID int64, formName, text string, files []*mautrix.ReqUploadMedia) {
	if issueID == 0 {
		return
	}

	log := zerolog.Ctx(ctx).With().Str("form", formName).Int64("issue_id", issueID).Logger()
	var attachments []*redmine.UploadRequest
	if len(files) > 0 {
		attachments = []*redmine.UploadRequest{}
		for _, media := range files {
			attachments = append(attachments, &redmine.UploadRequest{
				Stream: bytes.NewReader(media.ContentBytes),
				Path:   media.FileName,
			})
		}
	}

	statusID := h.rdm.StatusToID(redmine.WaitingForOperator)
	if err := h.rdm.UpdateIssue(issueID, statusID, text, attachments...); err != nil {
		log.Error().Err(err).Msg("cannot update redmine issue")
	}
}

func (h *Handler) redirect(ctx context.Context, target string, vars map[string]string) string {
	log := zerolog.Ctx(ctx)
	span := utils.StartSpan(ctx, "sub.redirect")
	defer span.Finish()

	var html bytes.Buffer
	var targetBytes bytes.Buffer
	targetTpl, err := template.New("target").Parse(target)
	if err != nil {
		log.Error().Err(err).Msg("cannot parse redirect url template")
	}
	if targetTpl != nil {
		err = targetTpl.Execute(&targetBytes, vars)
		if err != nil {
			log.Error().Err(err).Msg("cannot execute redirect url template")
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
		log.Error().Err(err).Msg("cannot execute redirect template")
	}

	return html.String()
}

// generate text and files
func (h *Handler) generate(ctx context.Context, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	span := utils.StartSpan(ctx, "sub.generate")
	defer span.Finish()

	v := h.vs[form.Name]
	medias := []*mautrix.ReqUploadMedia{}
	text, rmedias := h.ext["root"].Execute(span.Context(), v, form, data)

	for _, extension := range form.Extensions {
		if extension == "" {
			continue
		}
		e, ok := h.ext[extension]
		if !ok || e == nil {
			continue
		}

		etext, emedias := e.Execute(span.Context(), v, form, data)
		text += etext
		medias = append(medias, emedias...)
	}
	medias = append(medias, rmedias...) // add submission.md at the end

	return text, medias
}

// extValidate validates submission with extensions
func (h *Handler) extValidate(ctx context.Context, form *config.Form, data map[string]string) error {
	span := utils.StartSpan(ctx, "sub.generate")
	defer span.Finish()

	v := h.vs[form.Name]
	for _, extension := range form.Extensions {
		if extension == "" {
			continue
		}
		e, ok := h.ext[extension]
		if !ok || e == nil {
			continue
		}

		if err := e.Validate(span.Context(), v, form, data); err != nil {
			return fmt.Errorf("%s: %w", extension, err)
		}
	}
	return nil
}
