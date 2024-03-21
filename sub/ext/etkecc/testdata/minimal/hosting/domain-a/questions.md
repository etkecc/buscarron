price: $15

```yaml
SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

```


___

```yaml
set -euxo pipefail
SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"example.com\",\"server_type\":\"cpx11\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}")
SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')
SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')
SERVER_IP4_ID=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.id')
SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')
curl -X "POST" "https://api.hetzner.cloud/v1/primary_ips/$SERVER_IP4_ID/actions/change_dns_ptr" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{ \"ip\": \"$SERVER_IP4\", \"dns_ptr\": \"matrix.example.com\" }"
curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"
echo -e "---\nHello,\n\nWe've received your payment and have prepared a server for you. Its IP addresses are:\n\n- IPv4: $SERVER_IP4\n- IPv6: $SERVER_IP6\n\nPlease, add the following DNS entries (ensure that the CloudFlare proxy is disabled, as it's known to cause issues with Matrix Federation):\n\n- matrix	A record	$SERVER_IP4\n- matrix	AAAA record	$SERVER_IP6\n- matrix	TXT record	v=spf1 ip4:$SERVER_IP4 ip6:$SERVER_IP6 -all\n- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;\n\nPlease let us know when you're ready with the DNS and delegation redirects (you could just add the @ record pointing to the matrix server IP instead) configuration, so we can proceed with your server's setup.\n\nRegards\n"
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



