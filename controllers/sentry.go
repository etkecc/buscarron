package controllers

import (
	"fmt"

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
			ctx := c.Request().Context()
			hub := sentry.GetHubFromContext(ctx)
			if hub == nil {
				hub = sentry.CurrentHub().Clone()
				ctx = sentry.SetHubOnContext(ctx, hub)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			options := []sentry.SpanOption{
				sentry.WithOpName("http.server"),
				sentry.ContinueFromRequest(c.Request()),
				sentry.WithTransactionSource(sentry.SourceURL),
			}

			transaction := sentry.StartTransaction(c.Request().Context(),
				fmt.Sprintf("%s %s", c.Request().Method, c.Path()),
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
