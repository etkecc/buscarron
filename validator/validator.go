package validator

import (
	"net"
	"regexp"
	"strings"

	"gitlab.com/etke.cc/buscarron/logger"
)

// V is a validator implementation
type V struct {
	hosts  []string
	emails []string
	log    *logger.Logger
}

// based on W3C email regex, ref: https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var (
	emailRegex  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	domainRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$`)
)

// New Validator
func New(spamHosts []string, spamEmails []string, loglevel string) *V {
	return &V{
		hosts:  spamHosts,
		emails: spamEmails,
		log:    logger.New("v.", loglevel),
	}
}

// Domain checks if domain is valid
func (v *V) Domain(domain string) bool {
	// edge case: domain may be optional
	if domain == "" {
		return true
	}

	if len(domain) < 4 || len(domain) > 77 {
		v.log.Info("domain %s invalid, reason: length", domain)
		return false
	}

	if !domainRegex.MatchString(domain) {
		v.log.Info("domain %s invalid, reason: regexp", domain)
		return false
	}

	if !v.NS(domain) {
		v.log.Info("domain %s invalid, reason: nslookup", domain)
		return false
	}

	return true
}

// GetBase returns base domain/host of the provided domain
func (v *V) GetBase(domain string) string {
	parts := strings.Split(domain, ".")
	size := len(parts)
	// If domain contains only 2 parts (or less) - consider it without subdomains
	if size <= 2 {
		return domain
	}

	return parts[size-2] + "." + parts[size-1]
}

// Email checks if email is valid
func (v *V) Email(email string) bool {
	// edge case: email may be optional
	if email == "" {
		return true
	}

	// email cannot too short and too big
	if len(email) < 3 || len(email) > 254 {
		v.log.Info("email %s invalid, reason: length", email)
		return false
	}

	// check format
	if !emailRegex.MatchString(email) {
		v.log.Info("email %s invalid, reason: regexp", email)
		return false
	}

	if v.spam(email) {
		v.log.Info("email %s invalid, reason: spamlist", email)
		return false
	}

	return v.emailDomain(email)
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

// NS checks if host has at least one NS record
func (v *V) NS(host string) bool {
	if host == "" {
		return false
	}
	ns, err := net.LookupNS(host)
	if err != nil {
		v.log.Error("cannot get NS records of %s: %v", host, err)
		return false
	}

	return len(ns) > 0
}

// emailDomain checks if email domain or host is invalid
func (v *V) emailDomain(email string) bool {
	at := strings.LastIndex(email, "@")
	domain := email[at+1:]
	host := v.GetBase(domain)

	if v.spam(domain) {
		v.log.Info("email %s domain %s invalid, reason: spamlist", email, domain)
		return false
	}
	if v.spam(host) {
		v.log.Info("email %s host %s invalid, reason: spamlist", email, host)
		return false
	}

	if !v.MX(domain) && !v.MX(host) {
		v.log.Info("email %s domain/host %s invalid, reason: no MX", email, domain)
		return false
	}

	return true
}

// spam checks spam lists for the item
func (v *V) spam(item string) bool {
	for _, address := range v.emails {
		if address == item {
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
