package controllers

import (
	"fmt"

	"github.com/etkecc/buscarron/internal/utils"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// SentryTransaction is a middleware that creates a new transaction for each request.
func SentryTransaction() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/_health" {
				return next(c)
			}
			ctx := utils.NewContext(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))
			options := []sentry.SpanOption{
				sentry.WithOpName("http.server"),
				sentry.ContinueFromRequest(c.Request()),
				sentry.WithTransactionSource(sentry.SourceURL),
			}

			path := c.Path()
			if path == "" || path == "/" {
				path = c.Request().URL.Path
			}

			transaction := sentry.StartTransaction(c.Request().Context(),
				fmt.Sprintf("%s %s", c.Request().Method, path),
				options...,
			)
			defer transaction.Finish()

			c.SetRequest(c.Request().WithContext(transaction.Context()))

			if err := next(c); err != nil {
				transaction.Status = sentry.HTTPtoSpanStatus(c.Response().Status)
				return err
			}
			transaction.Status = sentry.SpanStatusOK
			return nil
		}
	}
}
