package etkecc

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.com/etke.cc/go/pricify"
	"golang.org/x/text/cases"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"

	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/buscarron/utils"
)

type order struct {
	orderedAt time.Time
	test      bool
	name      string
	data      map[string]string
	pd        *pricify.Data
	pm        EmailSender
	v         common.Validator
	c         cases.Caser

	domain    string
	subdomain bool
	followup  *event.MessageEventContent
	hosting   string
	smtp      map[string]string
	dkim      map[string]string
	price     int

	txt   strings.Builder
	pass  map[string]string
	files []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

// execute converts order to special message and files
func (o *order) execute(ctx context.Context) (string, []*mautrix.ReqUploadMedia) {
	o.preprocess(ctx)
	o.txt.WriteString("price: $" + strconv.Itoa(o.price) + "\n\n")

	questions, countQ := o.generateQuestions(ctx)
	delegation := o.generateDelegationInstructions(ctx)
	dns := o.generateDNSInstructions(ctx)
	hosts := o.generateHosts()

	if countQ > 0 {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(questions)
		o.txt.WriteString("```\n\n")
		o.txt.WriteString("\n___\n\n")
	}

	if dns != "" {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(dns)
		o.txt.WriteString("```\n\n")
	}

	if delegation != "" {
		o.txt.WriteString("```yaml\n")
		o.txt.WriteString(delegation)
		o.txt.WriteString("```\n\n")
	}

	if hosts != "" {
		o.txt.WriteString("hosts:\n")
		o.txt.WriteString("```\n")
		o.txt.WriteString(hosts)
		o.txt.WriteString("\n")
		o.txt.WriteString("```\n\n")
	}
	o.txt.WriteString("\n\n")

	o.vars(ctx)
	o.generateOnboarding(ctx)
	o.generateFollowup(ctx, questions, delegation, dns, countQ)

	go o.toGP(ctx, hosts) //nolint:errcheck // no need to wait
	go o.sendFollowup(ctx)

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
func (o *order) preprocess(ctx context.Context) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.preprocess")
	defer span.Finish()

	for _, key := range preprocessFields {
		o.data[key] = strings.TrimSpace(strings.ToLower(o.data[key]))
	}
	o.data["homeserver"] = "synapse"
	o.data["matrix"] = "yes"

	o.hosting = o.getHostingSize()
	o.domain = o.v.GetBase(o.data["domain"])

	for suffix := range domains {
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
	o.preprocessPrice()

	o.password("matrix")
}

func (o *order) preprocessSMTP() {
	smtp := map[string]string{}
	if o.has("service-email") {
		smtp["host"] = "smtp.migadu.com"
		smtp["port"] = "587"
		smtp["login"] = "\"matrix@" + o.domain + "\""
		smtp["password"] = o.pwgen()
		smtp["email"] = smtp["login"]
	} else if o.has("smtp-relay") {
		smtp["host"] = o.get("smtp-relay-host")
		smtp["port"] = o.get("smtp-relay-port")
		smtp["login"] = o.get("smtp-relay-login")
		smtp["password"] = o.get("smtp-relay-password")
		smtp["email"] = o.get("smtp-relay-email")
	}
	o.smtp = smtp

	dkim := map[string]string{}
	dkim["record"], dkim["private"] = o.dkimgen()
	o.dkim = dkim
}

func (o *order) preprocessPrice() {
	if o.test {
		o.price = len(o.data)
		return
	}

	if o.pd == nil {
		return
	}

	o.price = o.pd.Calculate(o.data)
}
