package etkecc

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/go/pricify"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/format"

	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

type order struct {
	test bool
	name string
	data map[string]string
	pd   *pricify.Data
	pm   EmailSender
	v    common.Validator

	domain    string
	subdomain bool
	hosting   string
	smtp      map[string]string

	txt   strings.Builder
	eml   strings.Builder
	pass  map[string]string
	files []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

// execute converts order to special message and files
func (o *order) execute() (string, []*mautrix.ReqUploadMedia) {
	o.preprocess()

	questions, count := o.generateQuestions()
	dns, dnsInternal := o.generateDNSInstructions()
	hosts := o.generateHosts()

	o.txt.WriteString("```yaml\n")
	o.txt.WriteString(questions)
	o.txt.WriteString("```\n\n")
	o.txt.WriteString("\n___\n\n")

	if o.hosting != "" {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(o.generateHVPSCommand())
		o.txt.WriteString("```\n\n")
	}

	if o.hosting == "" || dnsInternal {
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
	if o.hosting == "" && !dnsInternal {
		o.eml.WriteString(dns)
	}
	o.eml.WriteString("\nPS: this is an automated email. Please, reply to it with answers to the questions above (if any). An operator (human) will proceed with your answers")

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
	return ok && val != "" && val != "no"
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
	o.data["matrix"] = "yes"

	o.hosting = o.getHostingSize()
	o.domain = o.v.GetBase(o.data["domain"])

	for suffix := range hDomains {
		if strings.HasSuffix(o.domain, suffix) {
			o.subdomain = true
			break
		}
	}

	o.data["serve_base_domain"] = "no"
	if !o.v.A(o.domain) && !o.v.CNAME(o.domain) {
		o.data["serve_base_domain"] = "yes"
	}

	if o.has("smtp-relay-password") {
		o.pass["smtp"] = o.get("smtp-relay-password")
	}
	o.preprocessSMTP()

	o.password("matrix")
}

func (o *order) preprocessSMTP() {
	smtp := map[string]string{}
	if o.has("service-email") {
		smtp["host"] = "smtp.migadu.com"
		smtp["port"] = "587"
		smtp["login"] = "\"matrix@{{ matrix_domain }}\""
		smtp["password"] = o.pwgen()
		smtp["email"] = smtp["login"]
	} else {
		smtp["host"] = o.get("smtp-relay-host")
		smtp["port"] = o.get("smtp-relay-port")
		smtp["login"] = o.get("smtp-relay-login")
		smtp["password"] = o.get("smtp-relay-password")
		smtp["email"] = o.get("smtp-relay-email")
	}
	o.smtp = smtp
}

func (o *order) sendmail() {
	if o.pm == nil || (reflect.ValueOf(o.pm).Kind() == reflect.Ptr && reflect.ValueOf(o.pm).IsNil()) {
		return
	}

	content := format.RenderMarkdown(o.eml.String(), true, true)
	req := &postmark.Email{
		To:       o.get("email"),
		Tag:      "confirmation",
		Subject:  "Matrix server on " + o.domain,
		TextBody: content.Body,
		HTMLBody: content.FormattedBody,
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
