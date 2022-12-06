```yaml
Hi there,
We got your order and have some questions before the setup.

SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

```


___

```yaml
curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"example.com\",\"server_type\":\"cx11\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}"
```

```yaml

DNS - please, add the following entries:
@	A record	server IP
matrix	A record	server IP
```

