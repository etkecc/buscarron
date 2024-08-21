

price: $7

[status page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947)

```yaml
Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/order/status/#ports-and-firewalls](https://etke.cc/order/status/#ports-and-firewalls)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
- matrix	TXT record	v=spf1 ip4:server IP -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
```

hosts:
```
example.com ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



