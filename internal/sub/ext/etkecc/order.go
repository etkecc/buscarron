package etkecc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/etkecc/go-pricify"
	"github.com/rs/zerolog"
	"golang.org/x/text/cases"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"

	"github.com/etkecc/buscarron/internal/sub/ext/common"
	"github.com/etkecc/buscarron/internal/utils"
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
	response  string
	hosting   string
	smtp      map[string]string
	dkim      map[string]string
	price     int

	txt    strings.Builder
	logins map[string]string
	pass   map[string]string
	files  []*mautrix.ReqUploadMedia
}

var preprocessFields = []string{"email", "domain", "username"}

func (o *order) logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().
		Str("domain", o.data["domain"]).
		Str("email", o.data["email"]).
		Logger()
}

// execute converts order to special message and files
func (o *order) execute(ctx context.Context) (htmlResponse, matrixMessage string, files []*mautrix.ReqUploadMedia) {
	log := o.logger(ctx)
	log.Info().Msg("starting order execution")
	o.preprocess(ctx)
	o.vars(ctx)
	o.generateOnboarding(ctx)

	o.txt.WriteString("\n\nprice: $" + strconv.Itoa(o.price) + "\n\n")

	h := sha256.New()
	h.Write([]byte(o.domain))
	id := hex.EncodeToString(h.Sum(nil))
	o.txt.WriteString("[status page](https://etke.cc/order/status/#" + id + ")\n\n")

	questions, countQ := o.generateQuestions(ctx)
	delegation := o.generateDelegationInstructions(ctx)
	dns := o.generateDNSInstructions(ctx)
	hosts := o.generateHosts(ctx)
	o.generateFollowup(ctx, questions, delegation, dns, countQ)

	go o.toGP(ctx, hosts) //nolint:errcheck // no need to wait
	go o.sendFollowup(ctx)

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

	log.Info().Msg("order has been executed")
	return o.response, o.txt.String(), o.files
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
	log := o.logger(ctx)
	log.Info().Msg("preprocessing order")
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

	if o.hosting != "" {
		delete(o.data, "ssh-host")
		delete(o.data, "ssh-user")
		delete(o.data, "ssh-password")
		delete(o.data, "ssh-port")
	}

	o.data["serve_base_domain"] = "no"
	if !o.v.A(o.domain) && !o.v.CNAME(o.domain) {
		o.data["serve_base_domain"] = "yes"
	}

	if o.has("smtp-relay-password") {
		o.pass["smtp"] = o.get("smtp-relay-password")
	}
	o.preprocessSMTP(span.Context())
	o.preprocessPrice(span.Context())
	o.preprocessS3(span.Context())
	o.preprocessSSH(span.Context())

	o.password("matrix")
}

func (o *order) preprocessSMTP(ctx context.Context) {
	log := o.logger(ctx)
	log.Info().Msg("preprocessing smtp")
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

func (o *order) preprocessS3(ctx context.Context) {
	log := o.logger(ctx)
	log.Info().Msg("preprocessing s3")
	// synapse needs https://host.url
	if o.has("synapse-s3-endpoint") {
		endpointURL, err := url.Parse(o.get("synapse-s3-endpoint"))
		if err == nil {
			endpointURL.Scheme = "https"
			o.data["synapse-s3-endpoint"] = endpointURL.String()
		}
	}

	// peertube and gotosocial need host.url
	if o.has("peertube-s3-endpoint") {
		endpointURL, err := url.Parse(o.get("peertube-s3-endpoint"))
		if err == nil {
			endpointURL.Scheme = ""
			o.data["peertube-s3-endpoint"] = strings.TrimPrefix(endpointURL.String(), "//")
		}
	}
	if o.has("gotosocial-s3-endpoint") {
		endpointURL, err := url.Parse(o.get("gotosocial-s3-endpoint"))
		if err == nil {
			endpointURL.Scheme = ""
			o.data["gotosocial-s3-endpoint"] = strings.TrimPrefix(endpointURL.String(), "//")
		}
	}
}

func (o *order) preprocessPrice(ctx context.Context) {
	log := o.logger(ctx)
	log.Info().Msg("preprocessing price")
	if o.test {
		o.price = len(o.data)
		log.Info().Msg("price has been preprocessed (test)")
		return
	}

	if o.pd == nil {
		log.Warn().Msg("price data is nil")
		return
	}

	o.price = o.pd.Calculate(o.data)
}

func (o *order) preprocessSSH(ctx context.Context) {
	log := o.logger(ctx)
	log.Info().Msg("preprocessing ssh")
	pub, priv := o.keygen()
	pub = strings.TrimSpace(pub) + " etke.cc"
	o.files = append(o.files, &mautrix.ReqUploadMedia{
		Content:       strings.NewReader(pub),
		ContentBytes:  []byte(pub),
		FileName:      "sshkey.pub",
		ContentType:   "text/plain",
		ContentLength: int64(len(pub)),
	},
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(priv),
			ContentBytes:  []byte(priv),
			FileName:      "sshkey.priv",
			ContentType:   "text/plain",
			ContentLength: int64(len(priv)),
		})
}
