package controllers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"slices"

	echobasicauth "github.com/etkecc/go-echo-basic-auth"
	"github.com/etkecc/go-kit"
	"github.com/etkecc/go-psd"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/etkecc/buscarron/internal/metrics"
	"github.com/etkecc/buscarron/internal/sub"
	"github.com/etkecc/buscarron/internal/sub/ext/common"
)

// FormHandler for web server
type FormHandler interface {
	GET(context.Context, string, *http.Request) (string, error)
	POST(context.Context, string, *http.Request) (string, error)
	SetSender(sender common.Sender)
}

type Config struct {
	FormHandler   FormHandler
	BanlistStatic []string
	BanlistSize   int
	FormRLsShared map[string]string
	FormRLs       map[string]string
	MetricsAuth   echobasicauth.Auth
	Validator     domainValidator
	PSD           *psd.Client
}

var formHandler FormHandler

// ConfigureRouter configures echo router
func ConfigureRouter(e *echo.Echo, cfg *Config) {
	formHandler = cfg.FormHandler
	banner := NewBanner(cfg.BanlistSize, cfg.BanlistStatic)
	validator := NewValidator(cfg.Validator, cfg.PSD)
	rl := NewRateLimiter(cfg.FormRLsShared, cfg.FormRLs)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return slices.Contains(donotban, c.Request().URL.Path)
		},
		Format:           `${custom} - - [${time_custom}] "${method} ${path} ${protocol}" ${status} ${bytes_out} "${referer}" "${user_agent}"` + "\n",
		CustomTimeFormat: "2/Jan/2006:15:04:05 -0700",
		CustomTagFunc: func(c echo.Context, w *bytes.Buffer) (int, error) {
			return w.Write([]byte(kit.AnonymizeIP(c.RealIP())))
		},
	}))
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(SentryTransaction())
	e.Use(middleware.Secure())
	e.Use(corsMiddleware())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderReferrerPolicy, "origin")
			return next(c)
		}
	})
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(true),
		echo.TrustLinkLocal(true),
		echo.TrustPrivateNet(true),
	)

	e.GET("/_health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/_domain", validator.DomainHander())
	e.GET("/metrics", echo.WrapHandler(&metrics.Handler{}), echobasicauth.NewMiddleware(&cfg.MetricsAuth))
	e.GET("/:name", func(c echo.Context) error {
		body, err := formHandler.GET(c.Request().Context(), c.Param("name"), c.Request())
		if errors.Is(err, sub.ErrNotFound) {
			return c.HTML(http.StatusNotFound, body)
		}
		if err != nil {
			return c.HTML(http.StatusInternalServerError, body)
		}

		return c.HTML(http.StatusOK, body)
	})
	e.POST("/:name", func(c echo.Context) error {
		log := zerolog.Ctx(c.Request().Context())
		body, err := formHandler.POST(c.Request().Context(), c.Param("name"), c.Request())
		if errors.Is(err, sub.ErrNotFound) || errors.Is(err, sub.ErrSpam) {
			log.Warn().Err(err).Str("name", c.Param("name")).Msg("submission has rejected")
			banner.Ban(c, err.Error())
			return c.HTML(http.StatusNotFound, body)
		}
		if err != nil {
			log.Error().Err(err).Str("name", c.Param("name")).Msg("submission has failed")
			return c.HTML(http.StatusInternalServerError, body)
		}

		return c.HTML(http.StatusOK, body)
	}, banner.Middleware(), rl.Middleware())
}

// SetFormHandlerSender sets sender for form handler
// it's a hack to avoid circular dependencies and allow setting matrix bot once it configured, while having HTTP server up & running, even if matrix part is down
func SetFormHandlerSender(sender common.Sender) {
	formHandler.SetSender(sender)
}
