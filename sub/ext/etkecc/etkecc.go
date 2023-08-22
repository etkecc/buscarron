package etkecc

import (
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc/pricify"
)

const pricifyDataURL = "https://gitlab.com/etke.cc/website/-/raw/svelte-order-form/svelte/src/lib/provider/dataConfig.json"

// Etkecc extension
type Etkecc struct {
	pm      EmailSender
	pricify *pricify.Data
	test    bool
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// New etke.cc extension
func New(pm EmailSender) *Etkecc {
	ext := &Etkecc{
		pm: pm,
	}
	ext.pricify, _ = pricify.New(pricifyDataURL) //nolint:errcheck // proof-of-concept
	return ext
}

// Execute extension
func (e *Etkecc) Execute(v common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	o := &order{
		name:  form.Name,
		data:  data,
		test:  e.test,
		v:     v,
		pd:    e.pricify,
		pm:    e.pm,
		pass:  map[string]string{},
		files: make([]*mautrix.ReqUploadMedia, 0, 3),
	}

	return o.execute()
}
