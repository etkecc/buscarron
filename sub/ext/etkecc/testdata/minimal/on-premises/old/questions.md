price: $5

```yaml
Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
```

hosts:
```
example.com ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



