```yaml
Hi there,
We got your order and have some questions before the setup.

```


___

```yaml
SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"example.com\",\"server_type\":\"cpx11\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}")
SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')
SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')
SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')
curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"
echo -e "---\nHello,\n\nWe've received your payment and have prepared a server for you. Its IP addresses are:\n\n- IPv4: $SERVER_IP4\n- IPv6: $SERVER_IP6\n\nPlease, add the following DNS entries:\n\n- @	A record	$SERVER_IP4\n- matrix	A record	$SERVER_IP4\n\nIf you care about IPv6, feel free to configure additional AAAA records in the steps mentioning A records above.\n\nPlease let us know when you're ready with the DNS configuration, so we can proceed with your server's setup.\n\nRegards\n"
```

hosts:
```
example.com ansible_host=1.2.3.4 ansible_user=user ansible_become_password=password ansible_port=222
```

questions: 0



**price**: $15/month
