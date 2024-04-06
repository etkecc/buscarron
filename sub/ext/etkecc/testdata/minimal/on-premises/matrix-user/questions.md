price: $12

```yaml

Please, add the following DNS entries:

- @	A record	1.2.3.4
- matrix	A record	1.2.3.4
- matrix	TXT record	v=spf1 ip4:1.2.3.4 -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA136YacxIZjvWDXDjqyyhBWhz5vnlGxRgUNR6qiUNQ62d8VulI7u3X/ulSExnguVqW0la910e+p47777V/6BUpmri9ph7MbieV4QnJw078FCMZV2yflC7r4Z02Y4OaWV9xEN6pFgaSNwd5va3TTCaLfDOOehQ820DHPP/YqUQhhnfF5A3QtI8dZdTZu3ljvAApwVG4GNguucKuRd7L8l3f/Ajbgq5bd9239vGS8Y2MSf7/RkoBdtzNTJZ370wc+lY7k9gLyxK40Fz8R0/DdYZRPgRquhA24AZjmf0IykTtqJKOrztY4XwdgHxuPEeXAcwTtpGfp1sMnBkp5gYJcUeuzAlKkh67oeeYxNUi5Lisn0v4+TAm+/EAbIve3jueOqbas9KpzsoAOPFUARZxZp+1k/6rZim7Bifc98Ef8xvNoJsI+Jd0ReZ2XB2WiWFxZepbYwF4/iC95CAV5BShF2+NJeZ+eLZlVmtscRT8axHLOLKcnGeaXx3zQm15DbIinsvAgMBAAE=
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=matrix ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



