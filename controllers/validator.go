package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

func (v *validator) DomainHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		domain := c.Request().URL.Query().Get("domain")
		if domain == "" || len(domain) < 11 {
			return c.NoContent(http.StatusNotFound)
		}
		domain = v.domain.GetBase(domain)

		if v.domain.DomainString(domain) && !v.domain.A(domain) {
			c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			return c.NoContent(http.StatusNoContent)
		}

		return c.NoContent(http.StatusConflict)
	}
}
