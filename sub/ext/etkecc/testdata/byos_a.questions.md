```yaml
Hi there,
We got your order and have some questions before the setup.

We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

DNS - please, add the following entries:
matrix	A record	server IP
etherpad	CNAME record	matrix.example.com
```

