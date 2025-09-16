package etkecc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/etkecc/go-kit"
	"github.com/etkecc/go-pricify"
	"github.com/etkecc/go-psd"
	"github.com/mattevans/postmark-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"maunium.net/go/mautrix"

	"github.com/etkecc/buscarron/internal/config"
	"github.com/etkecc/buscarron/internal/sub/ext/common"
	"github.com/etkecc/buscarron/internal/utils"
)

var psdc *psd.Client

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

// SetPSD sets PSD validator
func SetPSD(c *psd.Client) {
	psdc = c
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
func (e *Etkecc) Execute(ctx context.Context, v common.Validator, form *config.Form, data map[string]string) (htmlResponse, matrixMessage string, files []*mautrix.ReqUploadMedia) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.Execute")
	defer span.Finish()

	var p *pricify.Data
	var err error
	p, err = pricify.New()
	if err != nil && p == nil {
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
		logins:    map[string]string{},
		files:     make([]*mautrix.ReqUploadMedia, 0, 3),
	}

	return o.execute(span.Context())
}

// Validate submission
func (e *Etkecc) Validate(ctx context.Context, v common.Validator, _ *config.Form, data map[string]string) error {
	if psdc == nil {
		return nil
	}
	domain := strings.ToLower(strings.TrimSpace(data["domain"]))
	if domain == "" {
		return fmt.Errorf("domain is empty")
	}
	domain = v.GetBase(domain)
	targets, _ := psdc.GetWithContext(ctx, domain) //nolint:errcheck // that's ok
	if len(targets) > 0 {
		return fmt.Errorf("domain already exists")
	}

	if ip, ok := data["ssh-host"]; ok {
		if !kit.IsValidIP(ip) {
			return fmt.Errorf("invalid IP address")
		}

		targets, _ := psdc.GetWithContext(ctx, ip) //nolint:errcheck // that's ok
		if len(targets) > 0 {
			return fmt.Errorf("domain already exists")
		}
	}

	return nil
}

// PrivateSuffixes returns private suffixes
func PrivateSuffixes() []string {
	keys := make([]string, 0, len(domains))
	for k := range domains {
		keys = append(keys, k)
	}
	return keys
}
