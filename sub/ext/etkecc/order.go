package etkecc

import (
	"strings"

	"maunium.net/go/mautrix"
)

type order struct {
	test bool
	name string
	data map[string]string
	v    NetworkValidator

	txt   strings.Builder
	pass  map[string]string
	files []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

// execute converts order to special message and files
func (o *order) execute() (string, []*mautrix.ReqUploadMedia) {
	o.preprocess()

	o.txt.WriteString(o.generateQuestions())
	o.txt.WriteString("\n___\n\n")

	o.txt.WriteString("```yaml\n")
	o.txt.WriteString(o.generateDNS())
	o.txt.WriteString("```\n\n")

	o.generateVars()
	o.generateOnboarding()

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
