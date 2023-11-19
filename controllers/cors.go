package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func corsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if c.Request().Method == http.MethodOptions {
				c.Response().Header().Set("Access-Control-Max-Age", "86400")
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}
}
