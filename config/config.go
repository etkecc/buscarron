package config

import (
	"os"
	"strconv"
	"strings"

	"maunium.net/go/mautrix/id"
)

const prefix = "buscarron"

// New config
func New() *Config {
	cfg := &Config{
		Homeserver: env("homeserver", defaultConfig.Homeserver),
		Login:      env("login", defaultConfig.Login),
		Password:   env("password", defaultConfig.Password),
		Sentry:     env("sentry", defaultConfig.Sentry),
		LogLevel:   env("loglevel", defaultConfig.LogLevel),
		Port:       env("port", defaultConfig.Port),
		DB: DB{
			DSN:     env("db.dsn", defaultConfig.DB.DSN),
			Dialect: env("db.dialect", defaultConfig.DB.Dialect),
		},
		Ban: &Ban{
			Duration: envInt("ban.duration", defaultConfig.Ban.Duration),
			Size:     envInt("ban.size", defaultConfig.Ban.Size),
		},
		Spam: &Spam{
			Hosts:  envSlice("spam.hosts"),
			Emails: envSlice("spam.emails"),
		},
		Postmark: &Postmark{
			Token:   env("pm.token", ""),
			From:    env("pm.from", ""),
			ReplyTo: env("pm.replyto", ""),
		},
	}
	cfg.Forms = parseForms()

	return cfg
}

func env(shortkey string, defaultValue string) string {
	key := strings.ToUpper(prefix + "_" + strings.ReplaceAll(shortkey, ".", "_"))
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	return value
}

func envInt(shortkey string, defaultValue int) int {
	vString := env(shortkey, "")
	if vString == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(vString)
	if err != nil {
		return defaultValue
	}

	return value
}

func envSlice(shortkey string) []string {
	return strings.Split(env(shortkey, ""), " ")
}

func parseForms() map[string]*Form {
	list := envSlice("list")
	forms := make(map[string]*Form, len(list))
	for _, name := range list {
		form := &Form{
			RoomID:    id.RoomID(env(name+".room", "")),
			Name:      name,
			Redirect:  env(name+".redirect", ""),
			Ratelimit: env(name+".ratelimit", ""),
			Confirmation: Confirmation{
				Subject: env(name+".confirmation.subject", ""),
				Body:    env(name+".confirmation.body", ""),
			},
			Extensions: envSlice(name + ".extensions"),
		}
		forms[name] = form
	}
	return forms
}
