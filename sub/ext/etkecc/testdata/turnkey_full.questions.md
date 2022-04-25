```yaml
Hi there,
We got your order and have some questions before the setup.

I see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin

Honoroit: are you sure you want it? it's a helpdesk bot with e2e encryption support. Please, check the https://gitlab.com/etke.cc/honoroit and decide.

Telegram: please, go to https://my.telegram.org/apps and create a new app. Share the API ID and Hash with me

Uptime Kuma: are you sure you want it? It's a self-hosted monitoring tool, that is not integrated into the matrix, like 'opensource UptimeRobot'. Please, check the https://github.com/louislam/uptime-kuma and decide.

Radicale: are you sure you want it? It's a CalDAV/CardDAV server, that is not integrated into the matrix. Please, check the https://radicale.org/ and decide.

Miniflux: are you sure you want it? It's an RSS reader, not integrated into the matrix. Please, check the https://miniflux.app and decide.

Languagetool: are you sure you want it? It's an 'open-source Grammarly' server, requires ~30GB of disk space for n-grams, and is not integrated into the matrix. Please, check the https://languagetool.org and decide.

Soft-Serve: are you sure you want it? It's a ssh-git hosting, that is not integrated into the matrix. Please, check the https://github.com/charmbracelet/soft-serve and decide.

WireGuard and dnsmasq: are you sure you want them? WireGuard is a VPN (not integrated with matrix) and dnsmasq is a local DNS server installed with adblock lists (like pi-hole) that is integrated pretty well with WireGuard. Please, check the https://wireguard.com and decide. If you still want it, please, share with me a list of labels you want to assign to generated client keys (just to set filenames, so even '1,2,3...' is OK)

Etherpad (only with subscription or turnkey): are you sure you want it? It's a self-hosted collaborative editor, available to set up only with dimension (added a question about it, too). Keep in mind that you will get an etherpad anyway with the default integration manager, but it will be hosted by Element Inc. (developers of the Element client apps).

Dimension (only with subscription or turnkey): are you sure you want it? It's a self-hosted integration manager. You will get integration manager by default with any Element client app. Please check the https://github.com/turt2live/matrix-dimension and decide

Website (only with subscription or turnkey): to deploy a static website you have to point your base domain (the @ DNS entry) to the matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

SSO (only with subscription or turnkey): You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Please, send me the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal (only with subscription or turnkey): are you sure you want it? It's a push gateway, usable only for matrix client app developers, so you can't use it if you don't develop your mobile matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup (only with subscription or turnkey): please, provide the desired repository url (user@host:repo). We will generate ssh key and encryption passphrase on your server side. We will send you the public part of the generated ssh key. You will add that ssh key on your provider side.

email2matrix (only with subscription or turnkey): are you sure you want it? It's a one-way SMTP server to receive emails in a matrix room that is quite tricky to set up by you as it doesn't have a straightforward way to configure it, we will need to cooperate with you to do configuration both as matrix homeserver user (you, we don't have users on your homeserver thus don't have access to your data inside matrix) and system (us, because that tool's configuration available only in config files on the VM/VPS disk).

Jitsi (only with subscription or turnkey): are you sure you want it? You will get jitsi integration by default with public instance, the jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.

ma1sd (only with subscription or turnkey): are you sure you want it? It's deprecated software, previously used as stub - an identity server, unmaintained for a while, and it doesn't have any benefits in most cases (except if you want to add LDAP auth or Twilio phone number verification).

matrix-registration (only with subscription or turnkey): are you sure you want it? It's deprecated software - a workaround used to add invite-based registration to the matrix, because protocol didn't support it, but now you can use builtin invite tokens instead: https://matrix-org.github.io/synapse/latest/usage/administration/admin_api/registration_tokens.html

Miounne (only with subscription or turnkey): are you sure you want it? Miounne is deprecated software - a bridge between external services (like HTML/HTTP forms, matrix-registration, buymeacoffee, etc.) and matrix. Please, check the https://gitlab.com/etke.cc/miounne and decide. If you still want it, please, send me a configuration to apply (no, there is no 'default configuration'. No, there is no 'good configuration'. No, we don't provide configuration templates. It's completely up to you)

SSH: please, share with me YOUR public ssh key and YOUR public static IP(-s) to get ssh root access to your server. We restrict ssh server access by default to the predefined list of IPs and ssh keys to limit the attack surface. Of course, if you don't want to have ssh access or want to allow connections from anywhere (insecure) - just say a word.

```


___

```yaml

DNS - please, add the following entries:
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

