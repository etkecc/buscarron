package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
)

type ctxkey int

const (
	ctxID   ctxkey = iota
	ctxName ctxkey = iota
)

type ctxMiddleware struct {
	iph *iphasher
}

func (c *ctxMiddleware) Handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := c.iph.GetHash(r)
		name := strings.ReplaceAll(r.URL.Path, "/", "")

		ctx := context.WithValue(r.Context(), ctxID, id)
		ctx = context.WithValue(ctx, ctxName, name)

		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("ID", id)
				scope.SetExtra("name", name)
			})
		}

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}
}
