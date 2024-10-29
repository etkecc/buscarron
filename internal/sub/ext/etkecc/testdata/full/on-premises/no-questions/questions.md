

price: $104

[status page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947)

```yaml

Please, add the following DNS entries:

- @	A record	1.2.3.4
- matrix	A record	1.2.3.4
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
- matrix	TXT record	v=spf1 ip4:1.2.3.4 -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
- matrix	MX record	0 matrix.example.com.
- postmoogle._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



