```yaml
Hi there,
We got your order and have some questions before the setup.

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin

Honoroit: are you sure you want it? It's a helpdesk bot with e2e encryption support. Please, check [gitlab.com/etke.cc/honoroit](https://gitlab.com/etke.cc/honoroit) before deciding.

Telegram: please, go to [https://my.telegram.org/apps](https://https://my.telegram.org/apps) and create a new app. Share the API ID and Hash with us

SMTP relay: please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).

Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

buscarron: are you sure you want it? It's a bot that receives web forms (HTML/HTTP POST) and send them to (encrypted) Matrix rooms. Please, check https://gitlab.com/etke.cc/buscarron before deciding.

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - [github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs](https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs). Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup: please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

Jitsi: are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.

SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

```


___

```yaml
SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" -d "{\"name\":\"higenjitsuteki.etke.host\",\"server_type\":\"cpx21\",\"image\":\"ubuntu-22.04\",\"firewalls\":[{\"firewall\":124003}],\"ssh_keys\":[\"first\",\"second\",\"third\"],\"location\":\"fsn1\"}")
SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')
SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')
SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')
curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" -H "Content-Type: application/json" -H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"
echo -e "---\nHello,\n\nWe've received your payment and have prepared a server for you. Its IP addresses are:\n\n- IPv4: $SERVER_IP4\n- IPv6: $SERVER_IP6\n"
```

```yaml
export HETZNER_SERVER_IP=SERVER_IP
export MIGADU_VERIFICATION=CODE
curl -X "POST" "https://dns.hetzner.com/api/v1/records/bulk" -H "Content-Type: application/json" -H "Auth-API-Token: $HETZNER_API_TOKEN" -d "{\"records\":[{\"name\":\"higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"10 aspmx1.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"20 aspmx2.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"autoconfig.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"autoconfig.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"_autodiscover._tcp.higenjitsuteki\",\"type\":\"SRV\",\"value\":\"0 1 443 autodiscover.migadu.com\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"key1._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key1.higenjitsuteki.etke.host._domainkey.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"key2._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key2.higenjitsuteki.etke.host._domainkey.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"key3._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key3.higenjitsuteki.etke.host._domainkey.migadu.com.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"_dmarc.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 include:spf.migadu.com -all\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"hosted-email-verify=$MIGADU_VERIFICATION\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"buscarron.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"cinny.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"element.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"etherpad.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"hydrogen.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"jitsi.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"ntfy.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"schildichat.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"stats.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"},{\"name\":\"sygnal.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host.\",\"zone_id\":\"enTDpM8y67STAZcQMpmqr7\"}]}"
```

questions: 5



**price**: $145/month
