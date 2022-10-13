package etkecc

import (
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

// Etkecc extension
type Etkecc struct {
	pm   EmailSender
	test bool
}

// EmailSender interface
type EmailSender interface {
	Send(*postmark.Email) error
}

// New etke.cc extension
func New(pm EmailSender) *Etkecc {
	return &Etkecc{pm: pm}
}

// Execute extension
func (e *Etkecc) Execute(v common.Validator, form *config.Form, data map[string]string) (string, []*mautrix.ReqUploadMedia) {
	return parseOrder(form.Name, data, v, e.pm, e.test)
}

func parseOrder(name string, data map[string]string, v common.Validator, pm EmailSender, test bool) (string, []*mautrix.ReqUploadMedia) {
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
