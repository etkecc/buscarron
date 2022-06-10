package config

import (
	"gitlab.com/etke.cc/go/env"
	"maunium.net/go/mautrix/id"
)

const prefix = "buscarron"

// New config
func New() *Config {
	env.SetPrefix(prefix)
	cfg := &Config{
		Homeserver: env.String("homeserver", defaultConfig.Homeserver),
		Login:      env.String("login", defaultConfig.Login),
		Password:   env.String("password", defaultConfig.Password),
		Sentry:     env.String("sentry", defaultConfig.Sentry),
		LogLevel:   env.String("loglevel", defaultConfig.LogLevel),
		Port:       env.String("port", defaultConfig.Port),
		DB: DB{
			DSN:     env.String("db.dsn", defaultConfig.DB.DSN),
			Dialect: env.String("db.dialect", defaultConfig.DB.Dialect),
		},
		Ban: &Ban{
			Duration: env.Int("ban.duration", defaultConfig.Ban.Duration),
			Size:     env.Int("ban.size", defaultConfig.Ban.Size),
		},
		Spam: &Spam{
			Hosts:  env.Slice("spam.hosts"),
			Emails: env.Slice("spam.emails"),
		},
		Postmark: &Postmark{
			Token:   env.String("pm.token", ""),
			From:    env.String("pm.from", ""),
			ReplyTo: env.String("pm.replyto", ""),
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
			RoomID:    id.RoomID(env.String(name+".room", "")),
			Name:      name,
			Redirect:  env.String(name+".redirect", ""),
			Ratelimit: env.String(name+".ratelimit", ""),
			Confirmation: Confirmation{
				Subject: env.String(name+".confirmation.subject", ""),
				Body:    env.String(name+".confirmation.body", ""),
			},
			Extensions: env.Slice(name + ".extensions"),
		}
		forms[name] = form
	}
	return forms
}
