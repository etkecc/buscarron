```yaml
Hi there,
We got your order and have some questions before the setup.

We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

SSH: please, send us YOUR public SSH key and YOUR public static IP(-s) to get SSH root access to your server. We restrict SSH server access by default to the predefined list of IPs and SSH keys to limit the attack surface. Of course, if you don't want to have SSH access or want to allow connections from anywhere (insecure) - just say the word.

```


___

```yaml

DNS - please, add the following entries:
matrix	A record	server IP
```

