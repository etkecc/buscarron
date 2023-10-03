price: $139

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
- buscarron	CNAME record	matrix.example.com
- cinny	CNAME record	matrix.example.com
- element	CNAME record	matrix.example.com
- etherpad	CNAME record	matrix.example.com
- social	CNAME record	matrix.example.com
- hydrogen	CNAME record	matrix.example.com
- jitsi	CNAME record	matrix.example.com
- miniflux	CNAME record	matrix.example.com
- ntfy	CNAME record	matrix.example.com
- radicale	CNAME record	matrix.example.com
- schildichat	CNAME record	matrix.example.com
- stats	CNAME record	matrix.example.com
- sygnal	CNAME record	matrix.example.com
- kuma	CNAME record	matrix.example.com
- vault	CNAME record	matrix.example.com
- matrix	MX record	matrix.example.com
- matrix	TXT record	v=spf1 ip4:server IP -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
```

Don't forget to create a new firewall called example.com, allow the following IPs to access SSH port (22/tcp): `1.2.3.4, 5.6.7.8` and attach it to the server. (if customer requested to lift IP restriction, attach the `open-ssh` firewall to the server)

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222
```



