package controllers

import (
	"net"
	"net/http"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
	"github.com/etkecc/go-psd"
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
	Email(string, string, ...net.IP) bool
}

type validator struct {
	v    domainValidator
	psdc *psd.Client
}

func NewValidator(v domainValidator, psdc *psd.Client) *validator {
	return &validator{v, psdc}
}

// Handler returns an echo handler function for domain and email validation
func (v *validator) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.QueryParam("domain") != "" {
			return v.domainHander(c)
		}
		if c.QueryParam("email") != "" {
			return v.emailHandler(c)
		}
		return c.JSON(http.StatusNotFound, map[string]any{
			"error_code": "invalid-request",
		})
	}
}

// domainHander is a handler for domain validation (GET requests, JSON response)
func (v *validator) domainHander(c echo.Context) error {
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

	domain = v.v.GetBase(domain)
	resp["base"] = domain
	if !v.v.DomainString(domain) {
		resp["error_code"] = errInvalidDomain
		utils.SetCachedValidation("domain-"+domain, false)
		return c.JSON(http.StatusNotFound, resp)
	}
	if v.v.A(domain) {
		resp["taken"] = true
		resp["error_code"] = errDomainTaken
		return c.JSON(http.StatusConflict, resp)
	}

	if targets, _ := v.psdc.Get(domain); len(targets) > 0 { //nolint:errcheck // that's fine
		resp["taken"] = true
		resp["error_code"] = errDomainTaken
		return c.JSON(http.StatusConflict, resp)
	}

	if isPunycode {
		resp["error_code"] = errPunycodeDomain
	}

	return c.JSON(http.StatusOK, resp)
}

// emailHandler is a handler for email validation (GET requests, no body response)
func (v *validator) emailHandler(c echo.Context) error {
	email := strings.ToLower(strings.TrimSpace(c.QueryParam("email")))
	if email == "" {
		return c.NoContent(http.StatusNoContent)
	}
	c.Logger().Infof("validating email %q", email)
	if valid, cached := utils.GetCachedValidation("email-" + email); cached {
		c.Logger().Infof("email %q is valid=%t (cached)", email, valid)
		return c.NoContent(http.StatusNoContent)
	}

	if !v.v.Email(email, "") {
		c.Logger().Infof("email %q is valid=%t (check)", email, false)
		utils.SetCachedValidation("email-"+email, false)
		return c.NoContent(http.StatusNoContent)
	}

	c.Logger().Infof("email %q is valid=%t (check)", email, true)
	utils.SetCachedValidation("email-"+email, true)
	return c.NoContent(http.StatusNoContent)
}
