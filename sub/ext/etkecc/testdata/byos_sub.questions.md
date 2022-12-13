```yaml
Hi there,
We got your order and have some questions before the setup.

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.
Add our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml
export HETZNER_SERVER_IP=SERVER_IP
curl -X "POST" "https://dns.hetzner.com/api/v1/records/bulk" -H "Content-Type: application/json" -H "Auth-API-Token: $HETZNER_API_TOKEN" -d "{\"records\":[{\"name\":\"higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"$HETZNER_ZONE_ID\"}]}"
```

