```yaml
Hi there,
We got your order and have some questions before the setup.

We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301):
* [example.com/.well-known/matrix/server](https://example.com/.well-known/matrix/server) -> [matrix.example.com/.well-known/matrix/server](https://matrix.example.com/.well-known/matrix/server)
* [example.com/.well-known/matrix/client](https://example.com/.well-known/matrix/client) -> [matrix.example.com/.well-known/matrix/client](https://matrix.example.com/.well-known/matrix/client)
To learn more about why these redirects are necessary and what the connection between the base domain (example.com) and the Matrix domain (matrix.example.com) is, read the following guide: [etke.cc/help/faq#why-do-i-need-well-known-redirects-on-the-base-domain](https://etke.cc/help/faq#why-do-i-need-well-known-redirects-on-the-base-domain)

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider)

```


___

```yaml

Please, add the following DNS entries:

- matrix	A record	server IP
- etherpad	CNAME record	matrix.example.com
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password
```

questions: 2



**price**: $7/month
