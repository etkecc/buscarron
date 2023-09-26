```yaml
Hi there,
We got your order and have some questions before the setup.

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin. [Full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)

Honoroit: are you sure you want it? It's a helpdesk bot with e2e encryption support. Please, check [gitlab.com/etke.cc/honoroit](https://gitlab.com/etke.cc/honoroit) before deciding.

Telegram: please, go to [https://my.telegram.org/apps](https://https://my.telegram.org/apps) and create a new app. Share the API ID and Hash with us

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider)

Prometheus+Grafana: are you sure you want it? Cloud providers usually provide a dashboard with server stats, so why not use that dashboard instead? A Prometheus+Grafana stack provides some internal Matrix stats (like count of events), but it's overkill if you just want to see server utilization.

Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

buscarron: are you sure you want it? It's a bot that receives web forms (HTML/HTTP POST) and send them to (encrypted) Matrix rooms. Please, check [gitlab.com/etke.cc/buscarron](https://gitlab.com/etke.cc/buscarron) before deciding.

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - [github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs](https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs). Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup: please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

Jitsi: are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/ssh.key](https://etke.cc/ssh.key)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
- buscarron	CNAME record	matrix.example.com
- cinny	CNAME record	matrix.example.com
- element	CNAME record	matrix.example.com
- etherpad	CNAME record	matrix.example.com
- hydrogen	CNAME record	matrix.example.com
- jitsi	CNAME record	matrix.example.com
- ntfy	CNAME record	matrix.example.com
- schildichat	CNAME record	matrix.example.com
- stats	CNAME record	matrix.example.com
- sygnal	CNAME record	matrix.example.com
- matrix	MX record	matrix.example.com
- matrix	TXT record	v=spf1 ip4:server IP -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
```

questions: 5



**price**: $30/month
