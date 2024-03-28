package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
	"gitlab.com/etke.cc/go/psd"
	"gitlab.com/etke.cc/go/validator/v2"
	"gitlab.com/etke.cc/linkpearl"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/bot"
	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/controllers"
	"gitlab.com/etke.cc/buscarron/mail"
	"gitlab.com/etke.cc/buscarron/sub"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
	"gitlab.com/etke.cc/buscarron/utils"
)

var (
	version = "development"
	mxb     *bot.Bot
	log     *zerolog.Logger
	e       *echo.Echo
)

func main() {
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

	initBot(cfg)
	initControllers(cfg)
	initShutdown()

	log.Debug().Msg("starting matrix bot...")
	go mxb.Start()
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("http server failed")
	}
}

func initBot(cfg *config.Config) {
	db, err := sql.Open(cfg.DB.Dialect, cfg.DB.DSN)
	if err != nil {
		log.Panic().Err(err).Msg("cannot initialize SQL database")
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
		log.Panic().Err(err).Msg("cannot initialize matrix bot")
	}
	mxb = bot.New(lp)
	log.Debug().Msg("bot has been created")
}

func initControllers(cfg *config.Config) {
	var rooms []id.RoomID
	srl := make(map[string]string)
	rls := make(map[string]string, len(cfg.Forms))
	frr := make(map[string]string, len(cfg.Forms))
	vs := make(map[string]common.Validator, len(cfg.Forms))
	for name, item := range cfg.Forms {
		rooms = append(rooms, item.RoomID)
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
	fh := sub.NewHandler(cfg.Forms, vs, pm, mxb)
	psdc := psd.NewClient(cfg.PSD.URL, cfg.PSD.Login, cfg.PSD.Password)
	etkecc.SetPSD(psdc)
	kfcfg := &controllers.KoFiConfig{
		VerificationToken: cfg.KoFiToken,
		Sender:            mxb,
		PaidMarker:        etkecc.MarkAsPaid,
		PSD:               psdc,
		Rooms:             rooms,
		Room:              id.RoomID(cfg.KoFiRoom),
	}

	srvv := validator.New(&validator.Config{Domain: validator.Domain{PrivateSuffixes: etkecc.PrivateSuffixes()}})
	ccfg := &controllers.Config{
		FormHandler:   fh,
		BanlistStatic: cfg.Ban.List,
		BanlistSize:   cfg.Ban.Size,
		FormRLsShared: srl,
		FormRLs:       rls,
		KoFiConfig:    kfcfg,
		MetricsAuth:   cfg.Metrics,
		Validator:     srvv,
	}
	e = echo.New()
	e.Logger = lecho.From(*log)
	controllers.ConfigureRouter(e, ccfg)
	log.Debug().Msg("web server has been configured")
}

func initShutdown() {
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for range listener {
			e.Shutdown(context.Background()) // nolint: errcheck
			mxb.Stop()
			sentry.Flush(5 * time.Second)

			os.Exit(0)
		}
	}()
}

func recovery() {
	defer sentry.Flush(2 * time.Second)
	err := recover()
	// no problem just shutdown
	if err == nil {
		return
	}

	sentry.CurrentHub().Recover(err)
	sentry.Flush(5 * time.Second)
}
