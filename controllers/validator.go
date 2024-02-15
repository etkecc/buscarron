package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/idna"
)

const (
	// errDomainRequired is returned when domain is not provided
	errDomainRequired = "domain-required"
	// errPunycodeDomain is returned when domain is IDNA (punycode)
	errPunycodeDomain = "punycode-domain"
	// errPunycodeFailed is returned when punycode conversion fails
	errPunycodeFailed = "punycode-failed"
	// errInvalidDomain is returned when domain is not valid (e.g. contains invalid characters)
	errInvalidDomain = "domain-invalid"
	// errDomainTaken is returned when domain is already taken (i.e., contains an A record)
	errDomainTaken = "domain-taken"
)

type domainValidator interface {
	A(string) bool
	DomainString(string) bool
	GetBase(string) string
}

type validator struct {
	domain domainValidator
}

func NewValidator(dv domainValidator) *validator {
	return &validator{dv}
}

// DomainValidatorHandler is a handler for domain validation (HEAD requests, no content)
func (v *validator) DomainValidatorHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		domain := c.Request().URL.Query().Get("domain")
		if domain == "" {
			return c.NoContent(http.StatusNotFound)
		}

		domain = v.domain.GetBase(domain)
		if !v.domain.DomainString(domain) {
			return c.NoContent(http.StatusNotFound)
		}

		if v.domain.A(domain) {
			return c.NoContent(http.StatusConflict)
		}

		c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		return c.NoContent(http.StatusNoContent)
	}
}

// DomainHander is a handler for domain validation (GET requests, JSON response)
func (v *validator) DomainHander() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := map[string]any{"base": "", "taken": false}
		domain := c.Request().URL.Query().Get("domain")
		if domain == "" {
			resp["error_code"] = errDomainRequired
			return c.JSON(http.StatusNotFound, resp)
		}

		var isPunycode bool
		asciiDomain, err := idna.ToASCII(domain)
		if err != nil {
			resp["error_code"] = errPunycodeFailed + " " + err.Error()
			return c.JSON(http.StatusNotFound, resp)
		}
		if domain != asciiDomain {
			isPunycode = true
			domain = asciiDomain
		}

		domain = v.domain.GetBase(domain)
		resp["base"] = domain
		if !v.domain.DomainString(domain) {
			resp["error_code"] = errInvalidDomain
			return c.JSON(http.StatusNotFound, resp)
		}
		if v.domain.A(domain) {
			resp["taken"] = true
			resp["error_code"] = errDomainTaken
			return c.JSON(http.StatusConflict, resp)
		}

		if isPunycode {
			resp["error_code"] = errPunycodeDomain
		}

		return c.JSON(http.StatusOK, resp)
	}
}
