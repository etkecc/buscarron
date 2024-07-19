package config

import (
	echobasicauth "gitlab.com/etke.cc/go/echo-basic-auth"
	"gitlab.com/etke.cc/go/env"
	"maunium.net/go/mautrix/id"
)

const prefix = "buscarron"

// New config
func New() *Config {
	env.SetPrefix(prefix)
	spamlist := migrateSpam(env.Slice("spam.emails"), env.Slice("spam.localparts"), env.Slice("spam.hosts"), env.Slice("spamlist"))
	cfg := &Config{
		Homeserver: env.String("homeserver"),
		Login:      env.String("login"),
		Password:   env.String("password"),
		Sentry:     env.String("sentry"),
		Healthchecks: Healthchecks{
			URL:  env.String("hc.url", defaultConfig.Healthchecks.URL),
			UUID: env.String("hc.uuid"),
		},
		Redmine: Redmine{
			Host:      env.String("redmine.host"),
			APIKey:    env.String("redmine.apikey"),
			ProjectID: env.String("redmine.project"),
			TrackerID: env.Int("redmine.trackerid"),
			StatusID:  env.Int("redmine.statusid"),
		},
		LogLevel: env.String("loglevel", defaultConfig.LogLevel),
		Port:     env.String("port", defaultConfig.Port),
		Metrics: echobasicauth.Auth{
			Login:    env.String("metrics.login"),
			Password: env.String("metrics.password"),
			IPs:      env.Slice("metrics.ips"),
		},
		DB: DB{
			DSN:     env.String("db.dsn", defaultConfig.DB.DSN),
			Dialect: env.String("db.dialect", defaultConfig.DB.Dialect),
		},
		PSD: PSD{
			URL:      env.String("psd.url"),
			Login:    env.String("psd.login"),
			Password: env.String("psd.password"),
		},
		Ban: &Ban{
			Size: env.Int("ban.size", defaultConfig.Ban.Size),
			List: env.Slice("ban.list"),
		},
		Spamlist: spamlist,
		Postmark: &Postmark{
			Token:   env.String("pm.token"),
			From:    env.String("pm.from"),
			ReplyTo: env.String("pm.replyto"),
		},
		SMTP: &SMTP{
			From:              env.String("smtp.from"),
			EnforceValidation: env.Bool("smtp.validation"),
		},
	}
	cfg.Forms = parseForms()

	return cfg
}

func parseForms() map[string]*Form {
	list := env.Slice("list")
	forms := make(map[string]*Form, len(list))
	for _, name := range list {
		form := &Form{
			RoomID:          id.RoomID(env.String(name + ".room")),
			Name:            name,
			Redirect:        env.String(name + ".redirect"),
			RejectRedirect:  env.String(name + ".redirect.reject"),
			Ratelimit:       env.String(name + ".ratelimit"),
			RatelimitShared: env.Bool(name + ".ratelimit.shared"),
			HasEmail:        env.Bool(name + ".hasemail"),
			HasDomain:       env.Bool(name + ".hasdomain"),
			Confirmation: Confirmation{
				Subject: env.String(name + ".confirmation.subject"),
				Body:    env.String(name + ".confirmation.body"),
			},
			Text:       env.String(name + ".text"),
			Extensions: env.Slice(name + ".extensions"),
		}
		forms[name] = form
	}
	return forms
}

func migrateSpam(emails, localparts, hosts, list []string) []string {
	uniq := map[string]struct{}{}
	for _, email := range emails {
		if email == "" {
			continue
		}
		uniq[email] = struct{}{}
	}

	for _, localpart := range localparts {
		if localpart == "" {
			continue
		}
		uniq[localpart+"@*"] = struct{}{}
	}

	for _, host := range hosts {
		if host == "" {
			continue
		}
		uniq["*@"+host] = struct{}{}
	}

	for _, item := range list {
		if item == "" {
			continue
		}
		uniq[item] = struct{}{}
	}

	spamlist := make([]string, 0, len(uniq))
	for item := range uniq {
		spamlist = append(spamlist, item)
	}

	return spamlist
}
