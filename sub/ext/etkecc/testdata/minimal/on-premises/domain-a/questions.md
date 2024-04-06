price: $6

```yaml
Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml

Please, add the following DNS entries (ensure that the CloudFlare proxy is disabled, as it's known to cause issues with Matrix Federation):

- matrix	A record	server IP
- matrix	TXT record	v=spf1 ip4:server IP -all
- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
- default._domainkey.matrix	TXT record	v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAtS+2q3Lgnilz4FGvPGURTpVZn5lTmDkLprN36qDIjBBfQCE9ngCoYR2P3rAWY04PGxsL8iFJNMvEqFWlR8PqZe+597K3JLcsRsIvNfdzcCCI6yoTelag9UvgqyRPADPmtQz3jilxcdPem38k9K6K+yDiNJGdKKmbTYKdVx1an9GfuWmEx4ZIQelUPGUeTnOZOUqCKZLC7NQsmwS4KGGkj03RH5EUGb7hagc63imgN1iAstB+ggiiS9eqTyKt2G9WC6U/0teTqiADM9Idpzmucaya/mDq/v5S9bysjD/7Iz2VOXBvKuGK3Wzb2eI/ZNEDArUh/XpWU5fKMYw+IJHObwBU92QAVegdtBVXaJiVsAD4SUmXDx7smzEr9mem1EcY0jq74mV8/ikAp58D7NHpHadBko0i/05JLht1RukdgZJCPf5sE6GJhUr1x4PU3BHM8ziBVFOVajxgukr/XtT1O7BBgHNXySG+mprVcMyt904IGPrT60Z+bhZHdmdgCiFnAgMBAAE=
```

```yaml

We see that you have something on your base domain.
**If** that's a domain registrar's (parking) page and/or you intend to serve base domain (example.com) from the matrix server, just add the `@` DNS record pointing to the server IP and tell us about that.
**If** that's your website and/or you intend to serve base domain from some other server, you should add the following HTTPS redirects (HTTP 301):
* [example.com/.well-known/matrix/server](https://example.com/.well-known/matrix/server) -> [matrix.example.com/.well-known/matrix/server](https://matrix.example.com/.well-known/matrix/server)
* [example.com/.well-known/matrix/client](https://example.com/.well-known/matrix/client) -> [matrix.example.com/.well-known/matrix/client](https://matrix.example.com/.well-known/matrix/client)
* [example.com/.well-known/matrix/support](https://example.com/.well-known/matrix/support) -> [matrix.example.com/.well-known/matrix/support](https://matrix.example.com/.well-known/matrix/support)
To learn more about why these redirects are necessary and what the connection between the base domain (example.com) and the Matrix domain (matrix.example.com) is, read the following guide: [etke.cc/help/faq#why-are-well-known-redirects-on-the-base-domain-important](https://etke.cc/help/faq#why-are-well-known-redirects-on-the-base-domain-important)

```

hosts:
```
example.com ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



