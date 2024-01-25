price: $5

```yaml
Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

```


___

```yaml
export SERVER_IP4=SERVER_IP
curl -X "POST" "https://dns.hetzner.com/api/v1/records/bulk" -H "Content-Type: application/json" -H "Auth-API-Token: $HETZNER_API_TOKEN" -d "{\"records\":[{\"name\":\"higenjitsuteki\",\"type\":\"A\",\"value\":\"$SERVER_IP4\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"A\",\"value\":\"$SERVER_IP4\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 ip4:$SERVER_IP4 -all\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"_dmarc.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"}]}"
```

hosts:
```
higenjitsuteki.onmatrix.chat ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



