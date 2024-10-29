

price: $58

[status page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947)

```yaml
Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin. [Full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)

Telegram: please, go to [https://my.telegram.org/apps](https://https://my.telegram.org/apps) and create a new app. Share the API ID and Hash with us

Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Supported generators: hugo, jekyll, plain html (no generator). Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - [github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs](https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs). Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup: please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 2vCPU, 2GB RAM.
Open the required ports ([etke.cc/order/status/#ports-and-firewalls](https://etke.cc/order/status/#ports-and-firewalls)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
- buscarron	CNAME record	matrix.example.com.
- cinny	CNAME record	matrix.example.com.
- element	CNAME record	matrix.example.com.
- etherpad	CNAME record	matrix.example.com.
- firezone	CNAME record	matrix.example.com.
- funkwhale	CNAME record	matrix.example.com.
- social	CNAME record	matrix.example.com.
- hydrogen	CNAME record	matrix.example.com.
- jitsi	CNAME record	matrix.example.com.
- languagetool	CNAME record	matrix.example.com.
- linkding	CNAME record	matrix.example.com.
- miniflux	CNAME record	matrix.example.com.
- ntfy	CNAME record	matrix.example.com.
- peertube	CNAME record	matrix.example.com.
- radicale	CNAME record	matrix.example.com.
- schildichat	CNAME record	matrix.example.com.
- stats	CNAME record	matrix.example.com.
- sygnal	CNAME record	matrix.example.com.
- kuma	CNAME record	matrix.example.com.
- vault	CNAME record	matrix.example.com.
- matrix	TXT record	v=spf1 ip4:server IP -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
- matrix	MX record	0 matrix.example.com.
- postmoogle._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
- @	MX record	10 aspmx1.migadu.com.
- @	MX record	20 aspmx2.migadu.com.
- @	TXT record	v=spf1 include:spf.migadu.com -all
- autoconfig	CNAME record	autoconfig.migadu.com.
- key1._domainkey	CNAME record	key1.example.com._domainkey.migadu.com.
- key2._domainkey	CNAME record	key2.example.com._domainkey.migadu.com.
- key3._domainkey	CNAME record	key3.example.com._domainkey.migadu.com.
- _dmarc	TXT record	v=DMARC1; p=quarantine;
- _autodiscover._tcp	SRV record	0 1 443 autodiscover.migadu.com
```

hosts:
```
example.com ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



