package validator

import (
	"net"
	"net/mail"
	"regexp"
	"strings"

	"gitlab.com/etke.cc/go/logger"
	"gitlab.com/etke.cc/go/trysmtp"
	"golang.org/x/net/publicsuffix"
)

// Validator interface
type Validator interface {
	Domain(string, bool) bool
	DomainString(string) bool
	Email(string) bool
	A(string) bool
	CNAME(string) bool
	MX(string) bool
	GetBase(domain string) string
}

// V is a validator implementation
type V struct {
	hosts       []string
	emails      []string
	localparts  []string
	from        string
	enforceSMTP bool
	log         *logger.Logger
}

// based on W3C email regex, ref: https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var (
	domainRegex   = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$`)
	privateSuffix = []string{".etke.host"}
)

// New Validator
func New(spamHosts []string, spamEmails []string, spamLocalparts []string, smtpFrom string, smtpEnforce bool, loglevel string) Validator {
	return &V{
		hosts:       spamHosts,
		emails:      spamEmails,
		localparts:  spamLocalparts,
		from:        smtpFrom,
		enforceSMTP: smtpEnforce,
		log:         logger.New("v.", loglevel),
	}
}

// Domain checks if domain is valid
func (v *V) Domain(domain string, must bool) bool {
	// edge case: domain may be optional
	if domain == "" && !must {
		return true
	}

	if !v.DomainString(domain) {
		return false
	}

	return true
}

// DomainString checks if domain string / value is valid using string checks like length and regexp
func (v *V) DomainString(domain string) bool {
	if len(domain) < 4 || len(domain) > 77 {
		v.log.Info("domain %s invalid, reason: length", domain)
		return false
	}

	if !domainRegex.MatchString(domain) {
		v.log.Info("domain %s invalid, reason: regexp", domain)
		return false
	}

	return true
}

// hasSuffix checks if domain has a suffix from public suffix list or from predefined suffix list
func (v *V) hasSuffix(domain string) bool {
	for _, suffix := range privateSuffix {
		if strings.HasSuffix(domain, suffix) {
			return true
		}
	}

	eTLD, _ := publicsuffix.PublicSuffix(domain)
	return strings.IndexByte(eTLD, '.') >= 0
}

// GetBase returns base domain/host of the provided domain
func (v *V) GetBase(domain string) string {
	// domain without subdomain "example.com" has parts: example com
	minSize := 2
	if v.hasSuffix(domain) {
		// domain with a certain TLDs contains 3 parts: example.co.uk -> example co uk
		minSize = 3
	}

	parts := strings.Split(domain, ".")
	size := len(parts)
	// If domain contains only 2 parts (or less) - consider it without subdomains
	if size <= minSize {
		return domain
	}

	// return domain without subdomain (sub.example.com -> example.com; sub.example.co.uk -> example.co.uk)
	return strings.Join(parts[size-minSize:], ".")
}

// Email checks if email is valid
func (v *V) Email(email string) bool {
	// edge case: email may be optional
	if email == "" {
		return true
	}

	length := len(email)
	// email cannot too short and too big
	if length < 3 || length > 254 {
		v.log.Info("email %s invalid, reason: length", email)
		return false
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		v.log.Info("email %s invalid, reason: %v", email, err)
		return false
	}

	if v.spam(email) {
		v.log.Info("email %s invalid, reason: spamlist", email)
		return false
	}

	localpart := email[:strings.LastIndex(email, "@")]
	if v.spam(localpart) {
		v.log.Info("email %s invalid, reason: spamlist", email)
		return false
	}

	if v.emailDomain(email) {
		return false
	}

	smtpCheck := !v.emailSMTP(email)
	if v.enforceSMTP {
		return smtpCheck
	}

	return true
}

// A checks if host has at least one A record
func (v *V) A(host string) bool {
	if host == "" {
		return false
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		v.log.Error("cannot get A records of %s: %v", host, err)
		return false
	}

	return len(ips) > 0
}

// CNAME checks if host has at least one CNAME record
func (v *V) CNAME(host string) bool {
	if host == "" {
		return false
	}
	cname, err := net.LookupCNAME(host)
	if err != nil {
		v.log.Error("cannot get CNAME records of %s: %v", host, err)
		return false
	}

	return cname != ""
}

// MX checks if host has at least one MX record
func (v *V) MX(host string) bool {
	if host == "" {
		return false
	}
	mxs, err := net.LookupMX(host)
	if err != nil {
		v.log.Error("cannot get MX records of %s: %v", host, err)
		return false
	}

	return len(mxs) > 0
}

// emailDomain checks if email domain or host is invalid
func (v *V) emailDomain(email string) bool {
	at := strings.LastIndex(email, "@")
	domain := email[at+1:]
	host := v.GetBase(domain)

	if v.spam(domain) {
		v.log.Info("email %s domain %s invalid, reason: spamlist", email, domain)
		return true
	}
	if v.spam(host) {
		v.log.Info("email %s host %s invalid, reason: spamlist", email, host)
		return true
	}

	if !v.MX(domain) && !v.MX(host) {
		v.log.Info("email %s domain/host %s invalid, reason: no MX", email, domain)
		return true
	}

	return false
}

func (v *V) emailSMTP(email string) bool {
	client, err := trysmtp.Connect(v.from, email)
	if err != nil {
		if strings.HasPrefix(err.Error(), "451") {
			v.log.Info("email %s may be invalid, reason: SMTP check (%v)", email, err)
			return false
		}

		v.log.Info("email %s invalid, reason: SMTP check (%v)", email, err)
		return true
	}
	defer client.Close()

	return false
}

// spam checks spam lists for the item
func (v *V) spam(item string) bool {
	for _, address := range v.emails {
		if address == item {
			return true
		}
	}

	for _, localpart := range v.localparts {
		if localpart == item {
			return true
		}
	}

	for _, host := range v.hosts {
		if host == item {
			return true
		}
	}

	return false
}
