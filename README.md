# Buscarron [![Matrix](https://img.shields.io/matrix/buscarron:etke.cc?logo=matrix&server_fqdn=matrix.org&style=for-the-badge)](https://matrix.to/#/#buscarron:etke.cc)[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/buscarron/badges/main/coverage.svg)](https://gitlab.com/etke.cc/buscarron/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/buscarron)](https://goreportcard.com/report/gitlab.com/etke.cc/buscarron) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/buscarron.svg)](https://pkg.go.dev/gitlab.com/etke.cc/buscarron)

> [more about that name](https://finalfantasy.fandom.com/wiki/Buscarron_Stacks)

Send any form (HTTP POST, HTML) to an encrypted matrix room, used as part of [etke.cc](https://etke.cc) service.

## Features

* End-to-End encryption
* Receive any HTTP POST request as form submission in a matrix room

## Configuration

env vars

### mandatory

* **BUSCARRON_HOMESERVER** - homeserver url, eg: `https://matrix.example.com`
* **BUSCARRON_LOGIN** - user login/localpart, eg: `buscarron`
* **BUSCARRON_PASSWORD** - user password
* **BUSCARRON_NAME_ROOM** - Matrix Room ID of the `name` form
* **BUSCARRON_NAME_REDIRECT** - URL to redirect after handling form `name` data submission, supports [go template](https://pkg.go.dev/text/template) vars from submission data (eg `email` field can be added to the redirect string as `{{ .email }}`)

### optional form configs

* **BUSCARRON_NAME_REDIRECT_REJECT** - URL to redirect after handling form `name` (doesn't support templates) for **rejected** submissions
* **BUSCARRON_NAME_RATELIMIT** - rate limit of the form, format: `<max requests>r/<interval:s,m>`, eg: `1r/s`, `54r/m`, `7r/h`
* **BUSCARRON_NAME_RATELIMIT_SHARED** - enable shared rate limit between forms with that option enabled
* **BUSCARRON_NAME_HASEMAIL** - if the form has an `email` field and you want to enforce email validation
* **BUSCARRON_NAME_HASDOMAIN** - if the form has a `domain` field and you want to enforce domain validation
* **BUSCARRON_NAME_EXTENSIONS** - space-separated list of the form extensions, allowed values: `confirmation`
* **BUSCARRON_NAME_TEXT** - custom form text, supports [go template](https://pkg.go.dev/text/template) vars from submission data (e.g., `email` field can be displayed as `{{ .email }}`)
* **BUSCARRON_NAME_CONFIRMATION_SUBJECT** - confrimation email subject, supports [go template](https://pkg.go.dev/text/template) vars from submission data (eg `email` field can be added to subject as `{{ .email }}`). Requires `confirmation` extension in the `BUSCARRON_NAME_EXTENSIONS` list and postmark configuration
* **BUSCARRON_NAME_CONFIRMATION_BODY** - confrimation email body, supports [go template](https://pkg.go.dev/text/template) vars from submission data (eg `email` field can be added to body as `{{ .email }}`). Requires `confirmation` extension in the `BUSCARRON_NAME_EXTENSIONS` list and postmark configuration


1. Add form name to **BUSCARRON_LIST**, eg: `export BUSCARRON_LIST="form1 form2"`)
2. Add form config, eg: `export BUSCARRON_FORM1_REDIRECT=https://http.cat/200`
3. Send POST request to the `http://127.0.0.1:8080/name`

### optional buscarron configs

* **BUSCARRON_SENTRY** - sentry DSN
* **BUSCARRON_HC_URL** - healthchecks.io URL, default: `https://hc-ping.com`
* **BUSCARRON_HC_UUID** - healthchecks.io check UUID
* **BUSCARRON_LOGLEVEL** - log level
* **BUSCARRON_DB_DSN** - database connection string
* **BUSCARRON_DB_DIALECT** - database dialect (postgres, sqlite3)
* **BUSCARRON_SPAMLIST** - list of spam emails with wildcards, eg: `export BUSCARRON_SPAMLIST=*@spammer.com annoy@gmail.com spammer@*`
* **BUSCARRON_BAN_SIZE** - jail size of banned users
* **BUSCARRON_BAN_LIST** - list of IP hashes / IDs for permanent ban
* **BUSCARRON_PM_TOKEN** - [Postmark](https://postmarkapp.com) server token
* **BUSCARRON_PM_FROM** - [Postmark](https://postmarkapp.com) sender signature
* **BUSCARRON_PM_REPLYTO** - reply-to email header
* **BUSCARRON_SMTP_FROM** - email address (from) for SMTP validation. Must be valid email on valid SMTP server, otherwise it will be rejected by other servers
* **BUSCARRON_SMTP_VALIDATION** - enforce SMTP validation
* **BUSCARRON_METRICS_LOGIN** - /metrics login
* **BUSCARRON_METRICS_PASSWORD** - /metrics password
* **BUSCARRON_METRICS_IPS** - /metrics allowed ips

### optional redmine configs

* **BUSCARRON_REDMINE_HOST** - redmine host, e.g. `https://redmine.example.com`
* **BUSCARRON_REDMINE_APIKEY** - redmine API key
* **BUSCARRON_REDMINE_PROJECT** - redmine project identifier, e.g. `internal-project`
* **BUSCARRON_REDMINE_TRACKERID** - redmine tracker ID, e.g. `1`
* **BUSCARRON_REDMINE_STATUSID** - redmine status ID, e.g. `1`

You can find default values in [config/defaults.go](config/defaults.go)

## Where to get

[docker registry](https://gitlab.com/etke.cc/buscarron/container_registry), [etke.cc](https://etke.cc)
