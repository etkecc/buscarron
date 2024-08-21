

price: $13

[status page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947)

```yaml

Please, add the following DNS entries:

- @	A record	1.2.3.4
- matrix	A record	1.2.3.4
- matrix	TXT record	v=spf1 ip4:1.2.3.4 -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=TODO
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=matrix ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



