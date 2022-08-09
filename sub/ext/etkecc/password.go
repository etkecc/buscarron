package etkecc

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"math/big"
	"strings"

	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

const (
	passlen = 64
	charset = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // a-z A-Z 0-9
)

var charsetlen = big.NewInt(57)

// pwgen is actual password generator
func (o *order) pwgen() string {
	if o.test {
		return "TODO"
	}
	var password strings.Builder

	for i := 0; i < passlen; i++ {
		// nolint // the configuration will be genered as template and must be modified manually after that, so even if password will not be generated that's not a problem
		index, _ := rand.Int(rand.Reader, charsetlen)
		password.WriteByte(charset[index.Int64()])
	}

	return password.String()
}

func (o *order) keygen() (string, string) {
	if o.test {
		return "ssh-todo TODO", "-----BEGIN OPENSSH PRIVATE KEY-----\nTODO\n-----END OPENSSH PRIVATE KEY-----"
	}
	publicBytes, privateBytes, _ := ed25519.GenerateKey(nil) //nolint:errcheck
	publicStruct, _ := ssh.NewPublicKey(publicBytes)         //nolint:errcheck

	pemblock := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privateBytes),
	}
	private := pem.EncodeToMemory(pemblock)
	public := ssh.MarshalAuthorizedKey(publicStruct)

	return string(public), string(private)
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
