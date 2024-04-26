package etkecc

import (
	"strconv"

	"gitlab.com/etke.cc/go/secgen"
)

const defaultPassLen = 64

// pwgen is actual password generator
func (o *order) pwgen(length ...int) string {
	passlen := defaultPassLen
	if len(length) > 0 {
		passlen = length[0]
	}

	if o.test {
		return "TODO" + strconv.Itoa(passlen)
	}

	return secgen.Password(passlen)
}

func (o *order) bytesgen(length ...int) string {
	passlen := defaultPassLen
	if len(length) > 0 {
		passlen = length[0]
	}

	if o.test {
		return "TODO" + strconv.Itoa(passlen)
	}

	return secgen.Base64Bytes(passlen)
}

func (o *order) keygen() (pub, priv string) {
	if o.test {
		return "ssh-todo TODO", "-----BEGIN OPENSSH PRIVATE KEY-----\nTODO\n-----END OPENSSH PRIVATE KEY-----"
	}
	pub, priv, _ = secgen.Keypair() //nolint:errcheck // error is always nil

	return pub, priv
}

func (o *order) dkimgen() (record, priv string) {
	if o.test {
		return "v=DKIM1; k=rsa; p=TODO", "-----BEGIN PRIVATE KEY-----\nTODO\n-----END PRIVATE KEY-----"
	}
	record, priv, _ = secgen.DKIM() //nolint:errcheck // error is always nil

	return record, priv
}

// password calls pwgen and saves result to internal map to export that password in multiple places (eg vars and onboarding)
func (o *order) password(name string) string {
	pass, ok := o.pass[name]
	if ok && pass != "" {
		return pass
	}

	pass = o.pwgen()
	o.pass[name] = pass
	return pass
}

// login saves login credentials to internal map to export that password in multiple places (eg vars and onboarding)
func (o *order) login(name string, value ...string) string {
	if len(value) > 0 {
		o.logins[name] = value[0]
		return value[0]
	}
	if v, ok := o.logins[name]; ok {
		return v
	}
	return o.get("username")
}
