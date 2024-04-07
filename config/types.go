package config

import (
	echobasicauth "gitlab.com/etke.cc/go/echo-basic-auth"
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
	// Metrics config
	Metrics echobasicauth.Auth
	// KoFiToken is Ko-Fi webhook verification token (do not use: will be excluded)
	KoFiToken string
	// KoFiRoom is Ko-Fi room when ko-fi webhook's target is not found
	KoFiRoom string
	// Forms map
	Forms map[string]*Form
	// Spamlist with wildcards
	Spamlist []string
	// Ban config
	Ban *Ban

	// PSD config
	PSD PSD

	// DB config
	DB DB

	// Sentry DSN
	Sentry string
	// Healthchecks.io config
	Healthchecks Healthchecks

	// Postmark config
	Postmark *Postmark

	// SMTP config
	SMTP *SMTP
}

// Healthchecks.io config
type Healthchecks struct {
	URL  string
	UUID string
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

// PSD config
type PSD struct {
	URL      string
	Login    string
	Password string
}

// Form config
type Form struct {
	// Name of the form
	Name string
	// Redirect is an url to redirect after submission
	Redirect string
	// RejectRedirect is an url to redirect after a rejected submission (both rate limit and spam)
	RejectRedirect string
	// RoomID to send submission
	RoomID id.RoomID
	// Ratelimit config
	Ratelimit string
	// RatelimitShared means that rate limit will be shared across other forms with that option enabled
	RatelimitShared bool
	// HasDomain enforces "domain" field check
	HasDomain bool
	// HasEmail enforces "email" field check
	HasEmail bool
	// Text template
	Text string
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
