package etkecc

import (
	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/buscarron/config"
	"maunium.net/go/mautrix"
)

// Etkecc extension
type Etkecc struct {
	v    NetworkValidator
	pm   EmailSender
	test bool
}

// NetworkValidator interface
type NetworkValidator interface {
	A(host string) bool
	CNAME(host string) bool
	GetBase(domain string) string
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// New etke.cc extension
func New(v NetworkValidator, pm EmailSender) *Etkecc {
	return &Etkecc{v: v, pm: pm}
}

// Execute extension
func (e *Etkecc) Execute(form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	return parseOrder(form.Name, data, e.v, e.pm, e.test)
}

func parseOrder(name string, data map[string]string, v NetworkValidator, pm EmailSender, test bool) (string, []*mautrix.ReqUploadMedia) {
	o := &order{
		name:  name,
		data:  data,
		test:  test,
		v:     v,
		pm:    pm,
		pass:  map[string]string{},
		files: make([]*mautrix.ReqUploadMedia, 0, 3),
	}

	return o.execute()
}
