package main

import (
	"database/sql"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/go/validator/v2"
	"gitlab.com/etke.cc/linkpearl"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/bot"
	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mail"
	"gitlab.com/etke.cc/buscarron/sub"
	"gitlab.com/etke.cc/buscarron/sub/ext/etkecc"
	"gitlab.com/etke.cc/buscarron/web"
)

var (
	version = "development"
	mxb     *bot.Bot
	srv     *web.Server
	log     *zerolog.Logger
)

func main() {
	cfg := config.New()

	initLog(cfg)
	defer recovery()

	log.Info().Msg("#############################")
	log.Info().Str("version", version).Msg("Buscarron")
	log.Info().Msg("Matrix: true")
	log.Info().Msg("HTTP: true")
	log.Info().Int("count", len(cfg.Forms)).Msg("Forms")
	log.Info().Msg("#############################")

	initBot(cfg)
	initSrv(cfg)
	initShutdown()

	log.Debug().Msg("starting matrix bot...")
	go mxb.Start()
	if err := srv.Start(); err != nil {
		log.Panic().Err(err).Msg("web server crashed")
	}
}

func initLog(cfg *config.Config) {
	loglevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		loglevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(loglevel)
	var w io.Writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
	sentryWriter, err := zlogsentry.New(cfg.Sentry)
	if err == nil {
		w = io.MultiWriter(sentryWriter, consoleWriter)
	} else {
		w = consoleWriter
	}
	logger := zerolog.New(w).With().Timestamp().Caller().Logger()
	log = &logger
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
	mxb = bot.New(lp, log)
	log.Debug().Msg("bot has been created")
}

func initSrv(cfg *config.Config) {
	var rooms []id.RoomID
	srl := make(map[string]string)
	rls := make(map[string]string, len(cfg.Forms))
	frr := make(map[string]string, len(cfg.Forms))
	vs := make(map[string]sub.Validator, len(cfg.Forms))
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
	pm := mail.New(cfg.Postmark.Token, cfg.Postmark.From, cfg.Postmark.ReplyTo, log)
	fh := sub.NewHandler(cfg.Forms, vs, pm, mxb, log)
	kfcfg := &web.KoFiConfig{
		VerificationToken: cfg.KoFiToken,
		Logger:            log,
		Sender:            mxb,
		Rooms:             rooms,
		Room:              id.RoomID(cfg.KoFiRoom),
	}

	srvv := validator.New(&validator.Config{Domain: validator.Domain{PrivateSuffixes: etkecc.PrivateSuffixes()}})
	srv = web.New(cfg.Port, srl, rls, frr, log, fh, srvv, kfcfg, cfg.Ban.Size, cfg.Ban.List)

	log.Debug().Msg("web server has been created")
}

func initShutdown() {
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for range listener {
			srv.Stop()
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
