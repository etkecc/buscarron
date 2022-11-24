package common

import "net"

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
