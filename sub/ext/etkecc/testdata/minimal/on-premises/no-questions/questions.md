price: $5

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
```

Don't forget to create a new firewall called example.com, allow the following IPs to access SSH port (22/tcp): `1.2.3.4, 5.6.7.8` and attach it to the server. (if customer requested to lift IP restriction, attach the `open-ssh` firewall to the server)

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222
```



