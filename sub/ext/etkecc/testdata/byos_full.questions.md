```yaml
Hi there,
We got your order and have some questions before the setup.

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin

Honoroit: are you sure you want it? it's a helpdesk bot with e2e encryption support. Please, check https://gitlab.com/etke.cc/honoroit before deciding.

Telegram: please, go to https://my.telegram.org/apps and create a new app. Share the API ID and Hash with me

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

Prometheus+Grafana: are you sure you want it? Cloud providers usually provide a dashboard with server stats, so why not use that dashboard instead? A Prometheus+Grafana stack provides some internal matrix stats (like count of events), but it's overkill if you just want to see server utilization.

Uptime Kuma: are you sure you want it? It's a self-hosted monitoring tool, that is not integrated into Matrix, like 'opensource UptimeRobot'. Please, check https://github.com/louislam/uptime-kuma before deciding.

Radicale: are you sure you want it? It's a CalDAV/CardDAV server, that is not integrated into Matrix. Please, check https://radicale.org/ before deciding.

Miniflux: are you sure you want it? It's an RSS reader, not integrated into Matrix. Please, check https://miniflux.app before deciding.

Languagetool: are you sure you want it? It's an 'open-source Grammarly' server, requires ~30GB of disk space for n-grams, and is not integrated into Matrix. Please, check https://languagetool.org before deciding.

Soft-Serve: are you sure you want it? It's a git hosting exposed over SSH, that is not integrated into Matrix. Please, check https://github.com/charmbracelet/soft-serve before deciding.

WireGuard and dnsmasq: are you sure you want them? WireGuard is a VPN (not integrated with Matrix) and dnsmasq is a local DNS server installed with adblock lists (like pi-hole) that is integrated pretty well with WireGuard. Please, check https://wireguard.com before deciding. If you still want it, please, send us a list of labels you want to assign to generated client keys (just to set filenames, so even '1,2,3...' is OK)

Etherpad (only with subscription or turnkey): are you sure you want it? It's a self-hosted collaborative editor, available for usage only with dimension (added a question about it, too). Keep in mind that you will get an etherpad anyway with the default integration manager, but it will be hosted by Element Inc. (developers of the Element client apps).

Dimension (only with subscription or turnkey): are you sure you want it? It's a self-hosted integration manager. You will get an integration manager by default with any Element client app. Please check https://github.com/turt2live/matrix-dimension before deciding

Website (only with subscription or turnkey): to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

SSO (only with subscription or turnkey): You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal (only with subscription or turnkey): are you sure you want it? It's a push gateway, usable only for matrix client app developers, so you can't use it if you don't develop your mobile matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup (only with subscription or turnkey): please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

email2matrix (only with subscription or turnkey): are you sure you want it? It's a one-way SMTP server to receive emails into a matrix room. It's quite tricky to set up by yourself, as it doesn't have a straightforward way to get it configured. We will need to cooperate with you to do configuration both as a matrix homeserver user (performed by you, because we don't have users on your homeserver and don't have access to your data inside matrix) and system (performed by us, because that tool's configuration is saved in config files on the VM/VPS disk).

Jitsi (only with subscription or turnkey): are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.

ma1sd (only with subscription or turnkey): are you sure you want it? It's deprecated software, previously used as stub - an identity server, unmaintained for a while, and it doesn't have any benefits in most cases (unless you want to add LDAP auth or Twilio phone number verification).

matrix-registration (only with subscription or turnkey): are you sure you want it? It's deprecated software - a workaround used to add invite-based registration support to matrix, because the protocol didn't support it, but now you can use the builtin invite tokens feature instead: https://matrix-org.github.io/synapse/latest/usage/administration/admin_api/registration_tokens.html

Miounne (only with subscription or turnkey): are you sure you want it? Miounne is deprecated software - a bridge between external services (like HTML/HTTP forms, matrix-registration, buymeacoffee, etc.) and matrix. Please, check https://gitlab.com/etke.cc/miounne before deciding. If you still want it, please, send me a configuration to apply (no, there is no 'default configuration'. No, there is no 'good configuration'. No, we don't provide configuration templates. It's completely up to you)

Server: please, create a VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

DNS - please, add the following entries:
@	A record	server IP
matrix	A record	server IP
buscarron	CNAME record	matrix.example.com
cinny	CNAME record	matrix.example.com
dimension	CNAME record	matrix.example.com
element	CNAME record	matrix.example.com
goneb	CNAME record	matrix.example.com
hydrogen	CNAME record	matrix.example.com
jitsi	CNAME record	matrix.example.com
kuma	CNAME record	matrix.example.com
languagetool	CNAME record	matrix.example.com
miniflux	CNAME record	matrix.example.com
miounne	CNAME record	matrix.example.com
radicale	CNAME record	matrix.example.com
stats	CNAME record	matrix.example.com
sygnal	CNAME record	matrix.example.com
```

