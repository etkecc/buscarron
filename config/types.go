package config

import (
	"sync"

	"maunium.net/go/mautrix/id"
)

// Config of Buscarron
type Config struct {
	// Homeserver url
	Homeserver string
	// Login is a MXID localpart (buscarron - OK, @buscarron:example.com - wrong)
	Login string
	// Password for login/password auth only
	Password string
	// LogLevel for logger
	LogLevel string
	// Port of the web server
	Port string
	// Forms map
	Forms map[string]*Form
	// Spam config
	Spam *Spam
	// Ban config
	Ban *Ban

	// DB config
	DB DB

	// Sentry DSN
	Sentry string

	// Postmark config
	Postmark *Postmark
}

// DB config
type DB struct {
	// DSN is a database connection string
	DSN string
	// Dialect of the db, allowed values: postgres, sqlite3
	Dialect string
}

// Spam config
type Spam struct {
	Hosts      []string
	Emails     []string
	Localparts []string
}

// Ban config
type Ban struct {
	Duration int
	Size     int
}

// Postmark config
type Postmark struct {
	Token   string
	From    string
	ReplyTo string
}

// Form config
type Form struct {
	sync.Mutex
	// Name of the form
	Name string
	// Redirect is an url to redirect after submission
	Redirect string
	// RoomID to send submission
	RoomID id.RoomID
	// Ratelimit config
	Ratelimit string
	// HasDomain if form has "domain" field
	HasDomain bool
	// Confirmation email config
	Confirmation Confirmation
	// Extensions list
	Extensions []string
}

// Confirmation email config
type Confirmation struct {
	// Subject of the confirmation email
	Subject string
	// Body of the confirmation email
	Body string
}
