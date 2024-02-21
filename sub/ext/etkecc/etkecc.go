package etkecc

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/go/pricify"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

// Etkecc extension
type Etkecc struct {
	pm      EmailSender
	pricify *pricify.Data
	now     func() time.Time
	test    bool
}

// EmailSender interface
type EmailSender interface {
	Send(context.Context, *postmark.Email) error
}

// New etke.cc extension
func New(pm EmailSender) *Etkecc {
	ext := &Etkecc{
		pm:  pm,
		now: time.Now,
	}
	ext.pricify, _ = pricify.New() //nolint:errcheck // proof-of-concept
	return ext
}

// Execute extension
func (e *Etkecc) Execute(ctx context.Context, v common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	span := sentry.StartSpan(ctx, "sub.ext.etkecc.Execute")
	defer span.Finish()

	var p *pricify.Data
	var err error
	p, err = pricify.New()
	if err != nil {
		p = e.pricify
	}

	o := &order{
		orderedAt: e.now().UTC(),
		name:      form.Name,
		data:      data,
		test:      e.test,
		v:         v,
		c:         cases.Title(language.English),
		pd:        p,
		pm:        e.pm,
		pass:      map[string]string{},
		files:     make([]*mautrix.ReqUploadMedia, 0, 3),
	}

	return o.execute(span.Context())
}

// PrivateSuffixes returns private suffixes
func PrivateSuffixes() []string {
	keys := make([]string, 0, len(hDomains))
	for k := range hDomains {
		keys = append(keys, k)
	}
	return keys
}
