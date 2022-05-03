```yaml
Hi there,
We got your order and have some questions before the setup.

We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

Server: please, create a VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH key (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

DNS - please, add the following entries:
matrix	A record	server IP
```

