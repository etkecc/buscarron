```yaml
Hi there,
We got your order and have some questions before the setup.

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

```


___

```yaml

Please, add the following DNS entries:

- @	A record	server IP
- matrix	A record	server IP
- etherpad	CNAME record	matrix.example.com
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password
```

questions: 1

