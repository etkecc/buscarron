package common

import (
	"bytes"
	"net"
	"text/template"
)

// Validator interface
type Validator interface {
	A(string) bool
	MX(string) bool
	CNAME(string) bool
	Email(string, ...net.IP) bool
	Domain(string) bool
	DomainString(string) bool
	GetBase(string) string
}

func ParseTemplate(tplString string, data map[string]string) (string, error) {
	var result bytes.Buffer
	tpl, err := template.New("template").Parse(tplString)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
