price: $12

```yaml

Please, add the following DNS entries:

- @	A record	1.2.3.4
- matrix	A record	1.2.3.4
- matrix	TXT record	v=spf1 ip4:1.2.3.4 -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAp3SW4C2QfBp9l2tJFKmjsyoGiY6rpb0lGaIiNCbt5wnnUIvpenSkzKgQhSMzz7pGodlfWsyx2OqOrIvAYZhwBUMpYs6ags9jaK+LMcVf2rdc1cd2Bd5JeQf8TrAEs8/gCQuXYWm5tX5njjw3sY55hTj1kEkxvk2d3vTULBGdO0eX0RmQemJBTxoFL4b0ZOSqU9II/SQy88hEZxhbmRMBGDaHyc/yxU9yOUMctgyOwis50rKg75rJPM2wreAaN+Lp6gwXG7aisnxIxoT5/W2emGlZfvQZpQzemO96Cf4SMBT0uJpqlCtGq6j9czalEMkVLhlr9Mk3Ks9oVMyQKd69RbXzTfEPsk7XWxSigOKTknDpTNulOQYCn/D3HnWaMqPFW9BXg/dKFE1SV0rGDMsXRzGb6s0gn7htZQt5NGtsOFcnpkiWn+KxgMEhePDBi0yeSso9q4D33PdBMuVASi0YjDqDEINc0lPxlR/2jD6C9Tkg2HpPoHHwT4ZDZ0NvnIVjAgMBAAE=
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



