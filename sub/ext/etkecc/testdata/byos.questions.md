```yaml
Hi there,
We got your order and have some questions before the setup.

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).

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



**price**: $2/month
