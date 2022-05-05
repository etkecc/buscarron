package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/etke.cc/linkpearl"
	lpcfg "gitlab.com/etke.cc/linkpearl/config"

	"gitlab.com/etke.cc/buscarron/bot"
	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/logger"
	"gitlab.com/etke.cc/buscarron/mail"
	"gitlab.com/etke.cc/buscarron/sub"
	"gitlab.com/etke.cc/buscarron/validator"
	"gitlab.com/etke.cc/buscarron/web"
)

var (
	version = "development"
	mxb     *bot.Bot
	srv     *web.Server
	log     *logger.Logger
)

func main() {
	cfg := config.New()
	log = logger.New("buscarron.", cfg.LogLevel)
	initSentry(cfg)
	defer recovery()

	log.Info("#############################")
	log.Info("Buscarron " + version)
	log.Info("Matrix: true")
	log.Info("HTTP: true")
	log.Info("Forms: %d", len(cfg.Forms))
	log.Info("#############################")

	initBot(cfg)
	initSrv(cfg)
	initShutdown()

	log.Debug("starting matrix bot...")
	go mxb.Start()
	if err := srv.Start(); err != nil {
		// nolint // log.Fatal calls panic, not exit
		log.Fatal("web server crashed: %v", err)
	}
}

func initSentry(cfg *config.Config) {
	env := version
	if env != "development" {
		env = "production"
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.Sentry,
		Release:     "buscarron@" + version,
		Environment: env,
	})
	if err != nil {
		log.Fatal("cannot initialize sentry: %v", err)
	}
}

func initBot(cfg *config.Config) {
	db, err := sql.Open(cfg.DB.Dialect, cfg.DB.DSN)
	if err != nil {
		log.Fatal("cannot initialize SQL database: %v", err)
	}

	mxlog := logger.New("matrix.", cfg.LogLevel)
	lp, err := linkpearl.New(&lpcfg.Config{
		Homeserver:   cfg.Homeserver,
		Login:        cfg.Login,
		Password:     cfg.Password,
		DB:           db,
		Dialect:      cfg.DB.Dialect,
		LPLogger:     mxlog,
		APILogger:    logger.New("api.", cfg.LogLevel),
		StoreLogger:  logger.New("store.", cfg.LogLevel),
		CryptoLogger: logger.New("olm.", cfg.LogLevel),
	})
	if err != nil {
		// nolint // Fatal = panic, not os.Exit()
		log.Fatal("cannot initialize matrix bot: %v", err)
	}
	mxb = bot.New(lp, mxlog)
	log.Debug("bot has been created")
}

func initSrv(cfg *config.Config) {
	rls := make(map[string]string, len(cfg.Forms))
	for name, item := range cfg.Forms {
		rls[name] = item.Ratelimit
	}
	v := validator.New(cfg.Spam.Hosts, cfg.Spam.Emails, cfg.LogLevel)
	pm := mail.New(cfg.Postmark.Token, cfg.Postmark.From, cfg.Postmark.ReplyTo, cfg.LogLevel)
	fh := sub.NewHandler(cfg.Forms, v, pm, mxb, cfg.LogLevel)
	srv = web.New(cfg.Port, rls, cfg.LogLevel, fh, v, cfg.Ban.Duration, cfg.Ban.Size)

	log.Debug("web server has been created")
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
