package controllers

import (
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/buscarron/metrics"
)

type banner struct {
	store *lru.Cache[string, struct{}]
}

var donotban = []string{"/favicon.ico", "/robots.txt", "/metrics", "/_domain", "/_health"}

// NewBanner creates banner
func NewBanner(size int, banlist []string) *banner {
	store, _ := lru.New[string, struct{}](size) //nolint:errcheck // only in case of size < 0
	for _, ip := range banlist {
		ip = strings.TrimSpace(ip)
		store.Add(ip, struct{}{})
		metrics.BanUser("persistent", "-")
	}

	return &banner{store}
}

// Ban by request data
func (b *banner) Ban(c echo.Context, reason string) {
	for _, exclude := range donotban {
		if c.Request().URL.Path == exclude {
			return
		}
	}
	ip := c.RealIP()
	b.store.Add(ip, struct{}{})
	metrics.BanUser(reason, c.Param("name"))

	zerolog.Ctx(c.Request().Context()).
		Info().
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("ip", ip).
		Msg("has been banned")
}

func (b *banner) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			if b.store.Contains(c.RealIP()) {
				metrics.BanRequest()
				zerolog.Ctx(c.Request().Context()).
					Debug().
					Str("method", c.Request().Method).
					Str("path", c.Request().URL.Path).
					Str("ip", ip).
					Msg("banned request attempt")
				return c.NoContent(http.StatusTooManyRequests)
			}
			return next(c)
		}
	}
}
