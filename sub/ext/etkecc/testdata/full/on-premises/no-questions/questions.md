price: $99

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
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1H0d1lDdw+bNAJ4JkbeDyWq/rp56Fq99Z8MMdrrc8Yaqoj2pFVw53up/TCMvtBirlEP/WOVS9ARK/i3ba5yk2t4sOe8FvGtJT9xPjcTWkGPAF5TW5yGs4W4NGebRfqAQddJLpTAT+Gfv3tmPtarJ3LMFhrPWuMLBw+h4FJ29T1YlPEKkshM6CPzo6caNeSezWqtlaIcgzR8GHNZSJ4rdBrey+TchZ90pBuM28jxKHeGsPoYptRIUi9RJmvzxDZfeALCj3bI4Yho88GpcxLdqBEDW1OOBLImQ57MJIYKsKTBUaavmksIq69vmzSZUVCdmg+eaZ5CdUW5+Sgzp3+bGVOIJ59hS9+ovIaHcIQ49+zafA9KJzMDzM6kPsPo+pwD7hJzgkOzWIYlqOIjoKE/Zz9blvnQJ2Y7wC7fjJj8jOTEej65nElaDidnqJrnOyC2ciV2oIMpTaSPNUCF9jAC2zjm2KC4t9OB7VZeFD+SwIa5mr9N/+jaZC7G8pqLhm1/DAgMBAAE=
- matrix	MX record	matrix.example.com.
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



