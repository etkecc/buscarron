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

func (o *order) keygen() (string, string) {
	if o.test {
		return "ssh-todo TODO", "-----BEGIN OPENSSH PRIVATE KEY-----\nTODO\n-----END OPENSSH PRIVATE KEY-----"
	}
	pub, priv, _ := secgen.Keypair() //nolint:errcheck

	return pub, priv
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
