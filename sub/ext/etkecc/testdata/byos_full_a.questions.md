Hi there,
We got your order and have some questions before the setup.

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin

Honoroit: are you sure you want it? it's a helpdesk bot with e2e encryption support. Please, check the https://gitlab.com/etke.cc/honoroit and decide.

Miounne: are you sure you want it? Miounne is deprecated software - a bridge between external services (like HTML/HTTP forms, matrix-registration, buymeacoffee, etc.) and matrix. Please, check the https://gitlab.com/etke.cc/miounne and decide. If you still want it, please, send me a configuration to apply (no, there is no 'default configuration'. No, there is no 'good configuration'. No, we don't provide configuration templates. It's completely up to you)

Telegram: please, go to https://my.telegram.org/apps and create a new app. Share the API ID and Hash with me

email2matrix: are you sure you want it? It's a one-way SMTP server to receive emails in a matrix room that is quite tricky to set up by you as it doesn't have a straightforward way to configure it, we will need to cooperate with you to do configuration both as matrix homeserver user (you, we don't have users on your homeserver thus don't have access to your data inside matrix) and system (us, because that tool's configuration available only in config files on the VM/VPS disk).

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send me an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

Jitsi: are you sure you want it? You will get jitsi integration by default with public instance, the jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.

ma1sd: are you sure you want it? It's deprecated software, previously used as stub - an identity server, unmaintained for a while, and it doesn't have any benefits in most cases (except if you want to add LDAP auth or Twilio phone number verification).

matrix-registration: are you sure you want it? It's deprecated software - a workaround used to add invite-based registration to the matrix, because protocol didn't support it, but now you can use builtin invite tokens instead: https://matrix-org.github.io/synapse/latest/usage/administration/admin_api/registration_tokens.html

Sygnal: are you sure you want it? It's a push gateway, usable only for matrix client app developers, so you can't use it if you don't develop your mobile matrix app

Prometheus+Grafana: are you sure you want it? Any cloud provider gives you a dashboard with server stats, why not use that dashboard? Prometheus+Grafana stack provides some internal matrix stats (like count of events), but it's overkill if you just want to see server utilization.

Uptime Kuma: are you sure you want it? It's a self-hosted monitoring tool, that is not integrated into the matrix, like 'opensource UptimeRobot'. Please, check the https://github.com/louislam/uptime-kuma and decide.

Radicale: are you sure you want it? It's a CalDAV/CardDAV server, that is not integrated into the matrix. Please, check the https://radicale.org/ and decide.

Miniflux: are you sure you want it? It's an RSS reader, not integrated into the matrix. Please, check the https://miniflux.app and decide.

Languagetool: are you sure you want it? It's an 'open-source Grammarly' server, requires ~30GB of disk space for n-grams, and is not integrated into the matrix. Please, check the https://languagetool.org and decide.

WireGuard and dnsmasq: are you sure you want them? WireGuard is a VPN (not integrated with matrix) and dnsmasq is a local DNS server installed with adblock lists (like pi-hole) that is integrated pretty well with WireGuard. Please, check the https://wireguard.com and decide. If you still want it, please, share with me a list of labels you want to assign to generated client keys (just to set filenames, so even '1,2,3...' is OK)

Etherpad: are you sure you want it? It's a self-hosted collaborative editor, available to set up only with dimension (added a question about it, too). Keep in mind that you will get an etherpad anyway with the default integration manager, but it will be hosted by Element Inc. (developers of the Element client apps).

Dimension: are you sure you want it? It's a self-hosted integration manager. You will get integration manager by default with any Element client app. Please check the https://github.com/turt2live/matrix-dimension and decide

Static website: are you sure you want it? To serve a static website, you must serve your base domain from the matrix server (set @ DNS record to the server IP address) and publish your **static** website into the public git repository. Only in that case it can be enabled and deployed automatically during the maintenance runs

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Please, send me the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Server: please, create a VPS with any Debian-based distro. Minimal comfortable configuration for a basic matrix server: 1vCPU, 2GB RAM.
Add my ssh key (https://etke.cc/ssh.key) to your server, share with me your server IP, the username (with permissions to call sudo), and password (if set).


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

