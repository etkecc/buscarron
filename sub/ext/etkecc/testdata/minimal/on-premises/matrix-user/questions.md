price: $5

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=matrix ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



