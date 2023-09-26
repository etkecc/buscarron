package etkecc

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"

	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc/pricify"
)

type order struct {
	test bool
	name string
	data map[string]string
	pd   *pricify.Data
	pm   EmailSender
	v    common.Validator

	txt   strings.Builder
	eml   strings.Builder
	pass  map[string]string
	files []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

// execute converts order to special message and files
func (o *order) execute() (string, []*mautrix.ReqUploadMedia) {
	o.preprocess()

	hostingSize := o.getHostingSize()
	questions, count := o.generateQuestions()
	dns, dnsInternal := o.generateDNSInstructions()
	hosts := o.generateHosts()

	o.txt.WriteString("```yaml\n")
	o.txt.WriteString(questions)
	o.txt.WriteString("```\n\n")
	o.txt.WriteString("\n___\n\n")

	if hostingSize != "" {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(o.generateHVPSCommand())
		o.txt.WriteString("```\n\n")
	}

	if hostingSize == "" || dnsInternal {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(dns)
		o.txt.WriteString("```\n\n")
	}

	if hosts != "" {
		o.txt.WriteString("hosts:\n")
		o.txt.WriteString("```\n")
		o.txt.WriteString(hosts)
		o.txt.WriteString("```\n\n")
	}

	o.txt.WriteString("questions: ")
	o.txt.WriteString(strconv.Itoa(count))
	o.txt.WriteString("\n\n")

	o.generateVars()
	o.generateOnboarding()

	o.eml.WriteString(questions)
	if hostingSize == "" && !dnsInternal {
		o.eml.WriteString(dns)
	}
	o.eml.WriteString("\n" + o.t("ps_automatic_email"))

	o.sendmail()
	o.pricify()

	return o.txt.String(), o.files
}

// get returns a value for the key from data store if exists, otherwise returns default value
func (o *order) get(key string) string {
	val, ok := o.data[key]
	if ok && val != "" {
		return val
	}

	val, ok = defaults[key]
	if ok {
		return val
	}
	return "TODO"
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

func (o *order) getHostingSize() string {
	if !o.has("turnkey") {
		return ""
	}

	value := o.get("turnkey")
	parts := strings.Split(value, "-")
	if len(parts) < 2 {
		return value // new approach
	}
	return parts[1] // old approach
}

// preprocess data
func (o *order) preprocess() {
	for _, key := range preprocessFields {
		o.data[key] = strings.TrimSpace(strings.ToLower(o.data[key]))
	}
	o.data["homeserver"] = "synapse"

	if o.name == "turnkey" || o.getHostingSize() != "" {
		o.data["type"] = "turnkey"
	}

	o.data["domain"] = o.v.GetBase(o.data["domain"])
	o.data["serve_base_domain"] = "no"
	if !o.v.A(o.data["domain"]) && !o.v.CNAME(o.data["domain"]) {
		o.data["serve_base_domain"] = "yes"
	}

	lang := o.get("lang")
	_, ok := i18n[lang]
	if !ok {
		o.data["lang"] = i18nDefault
	}

	o.password("matrix")
}

func (o *order) sendmail() {
	if o.pm == nil || (reflect.ValueOf(o.pm).Kind() == reflect.Ptr && reflect.ValueOf(o.pm).IsNil()) {
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

func (o *order) pricify() {
	if o.pd == nil {
		return
	}

	price := strconv.Itoa(o.pd.Calculate(o.data))
	o.txt.WriteString("\n\n**price**: $")
	o.txt.WriteString(price)
	o.txt.WriteString("/month\n")
}
