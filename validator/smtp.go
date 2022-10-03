package validator

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

var (
	// SMTPAddrs priority list
	SMTPAddrs = []string{":25", ":587", ":465"}
	// SMTPFrom address used to check SMTP connection
	SMTPFrom = "buscarron@ilydeen.org"
)

func (v *V) connectSMTP(email string) (*smtp.Client, error) {
	localname := strings.SplitN(SMTPFrom, "@", 2)[1]
	hostname := strings.SplitN(email, "@", 2)[1]

	v.log.Debug("MX lookup of %s", hostname)
	mxs, err := net.LookupMX(hostname)
	if err != nil {
		v.log.Error("cannot perform MX lookup: %v", err)
		return nil, err
	}

	for _, mx := range mxs {
		for _, addr := range SMTPAddrs {
			client := v.trySMTP(localname, strings.TrimSuffix(mx.Host, "."), addr)
			if client != nil {
				return client, nil
			}
		}
	}

	// If there are no MX records, according to https://datatracker.ietf.org/doc/html/rfc5321#section-5.1,
	// we're supposed to try talking directly to the host.
	if len(mxs) == 0 {
		for _, addr := range SMTPAddrs {
			client := v.trySMTP(localname, hostname, addr)
			if client != nil {
				return client, nil
			}
		}
	}

	return nil, fmt.Errorf("target SMTP server not found")
}

func (v *V) trySMTP(localname, mxhost, addr string) *smtp.Client {
	v.log.Debug("trying SMTP connection to %s%s", mxhost, addr)
	conn, err := smtp.Dial(mxhost + addr)
	if err != nil {
		v.log.Warn("cannot connect to the %s%s: %v", mxhost, addr, err)
		return nil
	}
	err = conn.Hello(localname)
	if err != nil {
		v.log.Warn("cannot call HELLO command of the %s%s: %v", mxhost, addr, err)
		return nil
	}
	if ok, _ := conn.Extension("STARTTLS"); ok {
		v.log.Debug("%s supports STARTTLS", mxhost)
		config := &tls.Config{ServerName: mxhost}
		err = conn.StartTLS(config)
		if err != nil {
			v.log.Warn("STARTTLS connection to the %s failed: %v", mxhost, err)
		}
	}

	return conn
}
