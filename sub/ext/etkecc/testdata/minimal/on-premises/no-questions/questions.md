price: $12

```yaml

Please, add the following DNS entries:

- @	A record	1.2.3.4
- matrix	A record	1.2.3.4
- matrix	TXT record	v=spf1 ip4:1.2.3.4 -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



