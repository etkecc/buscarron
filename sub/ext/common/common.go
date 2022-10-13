package common

// Validator interface
type Validator interface {
	A(string) bool
	MX(string) bool
	CNAME(string) bool
	Email(string) bool
	Domain(string) bool
	DomainString(string) bool
	GetBase(string) string
}
