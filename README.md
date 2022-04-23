# Buscarron [![Matrix](https://img.shields.io/matrix/buscarron:etke.cc?logo=matrix&style=for-the-badge)](https://matrix.to/#/#buscarron:etke.cc)[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/buscarron/badges/main/coverage.svg)](https://gitlab.com/etke.cc/buscarron/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/buscarron)](https://goreportcard.com/report/gitlab.com/etke.cc/buscarron) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/buscarron.svg)](https://pkg.go.dev/gitlab.com/etke.cc/buscarron)

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
* **BUSCARRON_NAME_REDIRECT** - URL to redirect after handling form `name` data
* **BUSCARRON_NAME_RATELIMIT** - rate limit of the form, format: `<max requests>r/<interval:s,m>`, eg: `1r/s` or `54r/m`

1. Add form name to **BUSCARRON_LIST**, eg: `export BUSCARRON_LIST="form1 form2"`)
2. Add form config, eg: `export BUSCARRON_FORM1_REDIRECT=https://http.cat/200`
3. Send POST request to the `http://127.0.0.1:8080/name`

### optional

* **BUSCARRON_SENTRY** - sentry DSN
* **BUSCARRON_LOGLEVEL** - log level
* **BUSCARRON_DB_DSN** - database connection string
* **BUSCARRON_DB_DIALECT** - database dialect (postgres, sqlite3)
* **BUSCARRON_SPAM_HOSTS** - list of spam domains, eg: `export BUSCARRON_SPAM_HOSTS="spammer.com notspammer.com"`
* **BUSCARRON_SPAM_EMAILS** - list of spam emails: eg: `export BUSCARRON_SPAM_EMAILS="annoy@gmail.com spammer@live.com"`

You can find default values in [config/defaults.go](config/defaults.go)

## Where to get

[docker registry](https://gitlab.com/etke.cc/buscarron/container_registry), [etke.cc](https://etke.cc)
