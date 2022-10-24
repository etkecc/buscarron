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
	// NoEncryption disables encryption
	NoEncryption bool
	// Port of the web server
	Port string
	// Forms map
	Forms map[string]*Form
	// Spamlist with wildcards
	Spamlist []string
	// Ban config
	Ban *Ban

	// DB config
	DB DB

	// Sentry DSN
	Sentry string

	// Postmark config
	Postmark *Postmark

	// SMTP config
	SMTP *SMTP
}

// DB config
type DB struct {
	// DSN is a database connection string
	DSN string
	// Dialect of the db, allowed values: postgres, sqlite3
	Dialect string
}

// Ban config
type Ban struct {
	Size int
	List []string
}

// Postmark config
type Postmark struct {
	Token   string
	From    string
	ReplyTo string
}

// SMTP config
type SMTP struct {
	From              string
	EnforceValidation bool
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
	// HasDomain enforces "domain" field check
	HasDomain bool
	// HasEmail enforces "email" field check
	HasEmail bool
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
