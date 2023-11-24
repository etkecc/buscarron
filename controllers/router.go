package controllers

import (
	"errors"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	echobasicauth "gitlab.com/etke.cc/go/echo-basic-auth"

	"gitlab.com/etke.cc/buscarron/metrics"
	"gitlab.com/etke.cc/buscarron/sub"
)

// FormHandler for web server
type FormHandler interface {
	GET(string, *http.Request) (string, error)
	POST(string, *http.Request) (string, error)
}

type Config struct {
	FormHandler   FormHandler
	BanlistStatic []string
	BanlistSize   int
	FormRLsShared map[string]string
	FormRLs       map[string]string
	MetricsAuth   echobasicauth.Auth
	KoFiConfig    *KoFiConfig
	Validator     domainValidator
	Logger        *zerolog.Logger
}

// ConfigureRouter configures echo router
func ConfigureRouter(e *echo.Echo, cfg *Config) {
	banner := NewBanner(cfg.BanlistSize, cfg.BanlistStatic, cfg.Logger)
	validator := NewValidator(cfg.Validator)
	kofi := NewKoFi(cfg.KoFiConfig)
	rl := NewRateLimiter(cfg.FormRLsShared, cfg.FormRLs, cfg.Logger)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return slices.Contains(donotban, c.Request().URL.Path)
		},
		Format:           `${remote_ip} - - [${time_custom}] "${method} ${path} ${protocol}" ${status} ${bytes_out} "${referer}" "${user_agent}"` + "\n",
		CustomTimeFormat: "2/Jan/2006:15:04:05 -0700",
	}))
	e.Use(middleware.Recover())
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

	// middleware.RateLimiter()
	e.GET("/_health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.HEAD("/_domain", validator.DomainHandler())
	e.POST("/_kofi", kofi.Handler())
	e.GET("/metrics", echo.WrapHandler(&metrics.Handler{}), echobasicauth.NewMiddleware(&cfg.MetricsAuth))
	e.GET("/:name", func(c echo.Context) error {
		body, err := cfg.FormHandler.GET(c.Param("name"), c.Request())
		if errors.Is(err, sub.ErrNotFound) || errors.Is(err, sub.ErrSpam) {
			return c.HTML(http.StatusNotFound, body)
		}
		if err != nil {
			return c.HTML(http.StatusInternalServerError, body)
		}

		return c.HTML(http.StatusOK, body)
	})
	e.POST("/:name", func(c echo.Context) error {
		body, err := cfg.FormHandler.POST(c.Param("name"), c.Request())
		if errors.Is(err, sub.ErrNotFound) || errors.Is(err, sub.ErrSpam) {
			cfg.Logger.Warn().Err(err).Str("name", c.Param("name")).Msg("submission has rejected")
			banner.Ban(c, err.Error())
			return c.HTML(http.StatusNotFound, body)
		}
		if err != nil {
			cfg.Logger.Error().Err(err).Str("name", c.Param("name")).Msg("submission has failed")
			return c.HTML(http.StatusInternalServerError, body)
		}

		return c.HTML(http.StatusOK, body)
	}, banner.Middleware(), rl.Middleware())
}
