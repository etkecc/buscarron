price: $188

```yaml
Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin. [Full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)

Telegram: please, go to [https://my.telegram.org/apps](https://https://my.telegram.org/apps) and create a new app. Share the API ID and Hash with us

Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Supported generators: hugo, jekyll, plain html (no generator). Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - [github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs](https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs). Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup: please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

```


___

```yaml
SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"higenjitsuteki.onmatrix.chat\",\"server_type\":\"cpx11\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}")
SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')
SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')
SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')
curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"
echo -e "---\nHello,\n\nWe've received your payment and have prepared a server for you. Its IP addresses are:\n\n- IPv4: $SERVER_IP4\n- IPv6: $SERVER_IP6\n"
```

```yaml
export HETZNER_SERVER_IP=SERVER_IP
export MIGADU_VERIFICATION=CODE
curl -X "POST" "https://dns.hetzner.com/api/v1/records/bulk" -H "Content-Type: application/json" -H "Auth-API-Token: $HETZNER_API_TOKEN" -d "{\"records\":[{\"name\":\"higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"10 aspmx1.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"20 aspmx2.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"autoconfig.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"autoconfig.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"_autodiscover._tcp.higenjitsuteki\",\"type\":\"SRV\",\"value\":\"0 1 443 autodiscover.migadu.com\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"key1._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key1.higenjitsuteki.onmatrix.chat._domainkey.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"key2._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key2.higenjitsuteki.onmatrix.chat._domainkey.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"key3._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key3.higenjitsuteki.onmatrix.chat._domainkey.migadu.com.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"_dmarc.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 include:spf.migadu.com -all\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"hosted-email-verify=$MIGADU_VERIFICATION\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"buscarron.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"cinny.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"element.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"etherpad.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"firezone.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"social.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"hydrogen.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"jitsi.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"linkding.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"miniflux.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"ntfy.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"radicale.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"schildichat.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"stats.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"sygnal.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"kuma.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"vault.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"MX\",\"value\":\"0 matrix.higenjitsuteki.onmatrix.chat.\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 ip4:$HETZNER_SERVER_IP -all\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"},{\"name\":\"_dmarc.matrix.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"zVNMf3dur7oHP8dcGETZs\"}]}"
```

hosts:
```
higenjitsuteki.onmatrix.chat ansible_host=TODO ordered_at=2021-01-01_00:00:00
```



