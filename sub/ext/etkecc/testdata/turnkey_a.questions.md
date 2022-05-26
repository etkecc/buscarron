```yaml
Hi there,
We got your order and have some questions before the setup.

We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

```


___

```yaml

DNS - please, add the following entries:
matrix	A record	server IP
```

