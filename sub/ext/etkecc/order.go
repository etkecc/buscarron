package etkecc

import (
	"strings"

	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"
)

type order struct {
	test bool
	name string
	data map[string]string
	pm   EmailSender
	v    NetworkValidator

	txt   strings.Builder
	eml   strings.Builder
	pass  map[string]string
	files []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

// execute converts order to special message and files
func (o *order) execute() (string, []*mautrix.ReqUploadMedia) {
	o.preprocess()

	questions := o.generateQuestions()
	dns := o.generateDNS()

	o.txt.WriteString("```yaml\n")
	o.txt.WriteString(questions)
	o.txt.WriteString("```\n\n")
	o.txt.WriteString("\n___\n\n")

	o.txt.WriteString("```yaml\n")
	o.txt.WriteString(dns)
	o.txt.WriteString("```\n\n")

	o.generateVars()
	o.generateOnboarding()

	o.eml.WriteString(questions)
	if o.get("type") != "turnkey" {
		o.eml.WriteString(dns)
	}
	o.eml.WriteString("\n" + o.t("ps_automatic_email"))

	o.sendmail()

	return o.txt.String(), o.files
}

// get returns a value for the key from data store if exists, otherwise returns default value
func (o *order) get(key string) string {
	val, ok := o.data[key]
	if ok && val != "" {
		return val
	}

	return defaults[key]
}

// has check if key exists in the data store
func (o *order) has(key string) bool {
	val, ok := o.data[key]
	return ok || val != ""
}

// t is a wrapper of the package's t with lang set
func (o *order) t(key string) string {
	return t(o.get("lang"), key)
}

// preprocess data
func (o *order) preprocess() {
	for _, key := range preprocessFields {
		o.data[key] = strings.TrimSpace(strings.ToLower(o.data[key]))
	}
	o.data["homeserver"] = "synapse"

	if o.name == "turnkey" {
		o.data["type"] = "turnkey"
	}

	o.data["domain"] = o.v.GetBase(o.data["domain"])
	o.data["serve_base_domain"] = "no"
	if !o.v.A(o.data["domain"]) && !o.v.CNAME(o.data["domain"]) {
		o.data["serve_base_domain"] = "yes"
	}

	if o.get("type") == "turnkey" {
		o.data["smtp-relay"] = "yes"
	}

	if o.has("etherpad") {
		o.data["dimension"] = "auto"
	}

	o.password("matrix")
}

func (o *order) sendmail() {
	if o.pm == nil {
		return
	}

	req := &postmark.Email{
		To:       o.get("email"),
		Tag:      "confirmation",
		Subject:  o.t("matrix_server_on") + " " + o.get("domain"),
		TextBody: o.eml.String(),
	}
	err := o.pm.Send(req)
	if err != nil {
		o.txt.WriteString("\n\n**confirmation email**: ❌\n")
		return
	}

	o.txt.WriteString("\n\n**confirmation email**: ✅\n")
}
