price: $15

```yaml
SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"example.com\",\"server_type\":\"cpx11\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}")
SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')
SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')
SERVER_IP4_ID=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.id')
SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')
curl -X "POST" "https://api.hetzner.cloud/v1/primary_ips/$SERVER_IP4_ID/actions/change_dns_ptr" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{ \"ip\": \"$SERVER_IP4\", \"dns_ptr\": \"matrix.example.com\" }"
curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"
curl -X "POST" "https://api.hetzner.cloud/v1/firewalls" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"example.com\",\"apply_to\":[{\"server\":{\"id\":$SERVER_ID},\"type\":\"server\"}],\"rules\":[{\"description\":\"SSH\",\"direction\":\"in\",\"port\":\"22\",\"protocol\":\"tcp\",\"source_ips\":[\"1.2.3.4/32\",\"5.6.7.8/32\"]}]}"
echo -e "---\nHello,\n\nWe've received your payment and have prepared a server for you. Its IP addresses are:\n\n- IPv4: $SERVER_IP4\n- IPv6: $SERVER_IP6\n\nPlease, add the following DNS entries:\n\n- @	A record	1.2.3.4\n- @	AAAA record	$SERVER_IP6\n- matrix	A record	1.2.3.4\n- matrix	AAAA record	$SERVER_IP6\n- matrix	TXT record	v=spf1 ip4:1.2.3.4 ip6:$SERVER_IP6 -all\n- _dmarc.matrix	TXT record	v=DMARC1; p=quarantine;\n\nPlease let us know when you're ready with the DNS configuration, so we can proceed with your server's setup.\n\nRegards\n"
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222 ordered_at=2021-01-01_00:00:00
```



