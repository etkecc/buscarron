package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/etkecc/go-healthchecks/v2"
	"github.com/etkecc/go-linkpearl"
	"github.com/etkecc/go-psd"
	"github.com/etkecc/go-redmine"
	"github.com/etkecc/go-validator/v2"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
	_ "modernc.org/sqlite"

	"github.com/etkecc/buscarron/internal/bot"
	"github.com/etkecc/buscarron/internal/config"
	"github.com/etkecc/buscarron/internal/controllers"
	"github.com/etkecc/buscarron/internal/mail"
	"github.com/etkecc/buscarron/internal/sub"
	"github.com/etkecc/buscarron/internal/sub/ext/common"
	"github.com/etkecc/buscarron/internal/sub/ext/etkecc"
	"github.com/etkecc/buscarron/internal/utils"
)

var (
	version = "development"
	mxb     *bot.Bot
	rdm     *redmine.Redmine
	log     *zerolog.Logger
	hc      *healthchecks.Client
	e       *echo.Echo
)

func main() {
	quit := make(chan struct{})
	cfg := config.New()
	utils.SetName("buscarron")
	utils.SetSentryDSN(cfg.Sentry)
	utils.SetLogLevel(cfg.LogLevel)
	log = zerolog.Ctx(utils.NewContext())

	defer recovery()

	log.Info().Msg("#############################")
	log.Info().Str("version", version).Msg("Buscarron")
	log.Info().Msg("Matrix: true")
	log.Info().Msg("HTTP: true")
	log.Info().Int("count", len(cfg.Forms)).Msg("Forms")
	log.Info().Msg("#############################")

	if cfg.Healthchecks.UUID != "" {
		log.Info().Str("url", cfg.Healthchecks.URL).Str("uuid", cfg.Healthchecks.UUID).Msg("Healthchecks enabled")
		hc = healthchecks.New(
			healthchecks.WithBaseURL(cfg.Healthchecks.URL),
			healthchecks.WithCheckUUID(cfg.Healthchecks.UUID),
		)
		hc.Start(strings.NewReader("buscarron is starting"))
		go hc.Auto(60 * time.Second)
	}

	var err error
	rdm, err = redmine.New(
		redmine.WithLog(log),
		redmine.WithHost(cfg.Redmine.Host),
		redmine.WithAPIKey(cfg.Redmine.APIKey),
		redmine.WithProjectIdentifier(cfg.Redmine.ProjectID),
		redmine.WithTrackerID(cfg.Redmine.TrackerID),
		redmine.WithWaitingForOperatorStatusID(cfg.Redmine.StatusID),
	)
	if err != nil {
		log.Warn().Err(err).Msg("cannot initialize redmine client")
	}
	if rdm.Enabled() {
		log.Info().Msg("redmine integration enabled")
	}

	initControllers(cfg, rdm)
	initShutdown(quit)

	for {
		if err := initBot(cfg); err != nil {
			log.Warn().Err(err).Msg("matrix bot startup failed, restarting in 10s...")
			hc.Fail(strings.NewReader(fmt.Sprintf("matrix bot startup failed: %+v, restarting in 10s...", err)))
			time.Sleep(10 * time.Second)
			continue
		}
		controllers.SetFormHandlerSender(mxb)
		break
	}

	// sometimes homeserver may be down for a few minutes, so we need to restart the service
	go func() {
		if !mxb.Enabled() {
			return
		}
		for {
			log.Info().Msg("starting matrix bot...")
			if err := mxb.Start(); err != nil {
				log.Warn().Err(err).Msg("matrix bot failed, restarting in 10s...")
				hc.Fail(strings.NewReader(fmt.Sprintf("matrix bot failed: %+v, restarting in 10s...", err)))
			}
			time.Sleep(10 * time.Second)
		}
	}()

	if err := e.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic().Err(err).Msg("http server failed")
	}
	<-quit
}

func initBot(cfg *config.Config) error {
	if cfg.DB.Dialect == "sqlite3" {
		cfg.DB.Dialect = "sqlite"
	}
	db, err := sql.Open(cfg.DB.Dialect, cfg.DB.DSN)
	if err != nil {
		return err
	}
	if cfg.DB.Dialect == "sqlite" {
		db.SetMaxOpenConns(1)
	}

	lp, err := linkpearl.New(&linkpearl.Config{
		Homeserver: cfg.Homeserver,
		Login:      cfg.Login,
		Password:   cfg.Password,
		DB:         db,
		Dialect:    cfg.DB.Dialect,
		Logger:     *log,
	})
	if err != nil {
		return err
	}
	mxb = bot.New(lp)
	return nil
}

func initControllers(cfg *config.Config, rdm *redmine.Redmine) {
	srl := make(map[string]string)
	rls := make(map[string]string, len(cfg.Forms))
	frr := make(map[string]string, len(cfg.Forms))
	vs := make(map[string]common.Validator, len(cfg.Forms))
	for name, item := range cfg.Forms {
		rls[name] = item.Ratelimit
		frr[name] = item.RejectRedirect
		if item.RatelimitShared {
			srl[name] = item.Ratelimit
		}

		vcfg := &validator.Config{
			Email: validator.Email{
				Enforce:  cfg.Forms[name].HasEmail,
				MX:       true,
				SMTP:     cfg.SMTP.EnforceValidation,
				Spamlist: cfg.Spamlist,
				From:     cfg.SMTP.From,
			},
			Domain: validator.Domain{
				Enforce:         cfg.Forms[name].HasDomain,
				PrivateSuffixes: etkecc.PrivateSuffixes(),
			},
		}
		v := validator.New(vcfg)
		vs[name] = v
	}
	pm := mail.New(cfg.Postmark.Token, cfg.Postmark.From, cfg.Postmark.ReplyTo)
	fh := sub.NewHandler(cfg.Forms, vs, pm, mxb, rdm)
	psdc := psd.NewClient(cfg.PSD.URL, cfg.PSD.Login, cfg.PSD.Password)
	etkecc.SetPSD(psdc)
	srvv := validator.New(&validator.Config{Domain: validator.Domain{PrivateSuffixes: etkecc.PrivateSuffixes()}})
	ccfg := &controllers.Config{
		FormHandler:   fh,
		BanlistStatic: cfg.Ban.List,
		BanlistSize:   cfg.Ban.Size,
		FormRLsShared: srl,
		FormRLs:       rls,
		MetricsAuth:   cfg.Metrics,
		Validator:     srvv,
		PSD:           psdc,
	}
	e = echo.New()
	e.Logger = lecho.From(*log)
	controllers.ConfigureRouter(e, ccfg)
	log.Debug().Msg("web server has been configured")
}

func initShutdown(quit chan struct{}) {
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-listener
		defer close(quit)

		e.Shutdown(context.Background()) //nolint:errcheck // we don't care about the error here
		mxb.Stop()
		sentry.Flush(5 * time.Second)
		if hc != nil {
			hc.Shutdown()
			hc.ExitStatus(0, strings.NewReader("buscarron is shutting down"))
		}

		os.Exit(0)
	}()
}

func recovery() {
	defer sentry.Flush(2 * time.Second)
	err := recover()
	// no problem just shutdown
	if err == nil {
		return
	}

	if hc != nil {
		hc.ExitStatus(1, strings.NewReader(fmt.Sprintf("panic: %+v", err)))
	}
	sentry.CurrentHub().Recover(err)
	sentry.Flush(5 * time.Second)
}
