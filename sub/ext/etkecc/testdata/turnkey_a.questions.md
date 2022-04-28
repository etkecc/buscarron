```yaml
Hi there,
We got your order and have some questions before the setup.

I see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

SSH: please, share with me YOUR public ssh key and YOUR public static IP(-s) to get ssh root access to your server. We restrict ssh server access by default to the predefined list of IPs and ssh keys to limit the attack surface. Of course, if you don't want to have ssh access or want to allow connections from anywhere (insecure) - just say a word.

```


___

```yaml

DNS - please, add the following entries:
matrix	A record	server IP
```

