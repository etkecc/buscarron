```yaml
Hallo zusammen,
Wir haben Ihre Bestellung erhalten und haben einige Fragen vor der Einrichtung.

Reminder bot: Was ist Ihre Zeitzone (IANA)? Wie America/Chicago, Asia/Seoul, oder Europe/Berlin

Honoroit: Sind Sie sicher, dass Sie es wollen? Es ist ein Helpdesk-Bot mit Ende-Zu-Ende-Verschlüsselungsunterstützung. Bitte prüfen Sie https://gitlab.com/etke.cc/honoroit, bevor Sie sich entscheiden.

Telegram: gehen Sie bitte auf https://my.telegram.org/apps und erstellen Sie eine neue App. Teilen Sie mir die API-ID und den Hash mit

SMTP relay: Bitte wählen Sie einen geeigneten E-Mail-Anbieter (große Anbieter wie Gmail oder Outlook sperren Sie für automatisierte E-Mails, daher müssen Sie einen Dienst finden, der den Versand von Verifizierungs-E-Mails erlaubt. Optional bieten wir einen solchen Dienst an). Bitte senden Sie uns einen SMTP Host, einen SMTP STARTTLS Port, ein SMTP Login, ein SMTP Passwort und eine SMTP E-Mail (normalerweise sind Login und E-Mail dasselbe, aber das hängt vom Anbieter ab).

Prometheus+Grafana: Sind Sie sicher, dass Sie das wollen? Cloud-Anbieter bieten in der Regel ein Dashboard mit Server-Statistiken, warum also nicht stattdessen dieses Dashboard verwenden? Ein Prometheus+Grafana-Stack bietet einige interne Matrix-Statistiken (z. B. die Anzahl der Ereignisse), ist aber zu viel des Guten, wenn Sie nur die Serverauslastung sehen wollen.

Dimension: Sind Sie sicher, dass Sie das wollen? Es ist ein selbst gehosteter Integrationsmanager. Sie erhalten standardmäßig einen Integrationsmanager mit jeder Element-Client-Anwendung. Bitte prüfen Sie https://github.com/turt2live/matrix-dimension bevor Sie sich entscheiden.

Website: Um eine statische Website einzurichten, müssen Sie Ihre Basisdomäne (den @ DNS eintrag) auf die IP-Adresse des Matrix-Servers verweisen, und der Quellcode der Website muss in einem öffentlichen git repo verfügbar sein. Sind Sie sicher, dass Sie das wollen? Wenn ja, geben Sie bitte die URL des Website-Repositorys, den Befehl (-s) zum Erstellen der Website und den Ordner an, in dem die Build-Distribution gespeichert ist (normalerweise public oder dist).

buscarron: Sind Sie sicher, dass Sie es wollen? Es ist ein Bot, der Webformulare (HTML/HTTP POST) empfängt und sie an (verschlüsselte) Matrixräume sendet. Bitte prüfen Sie https://gitlab.com/etke.cc/buscarron, bevor Sie sich entscheiden.

SSO: Sie haben nicht erwähnt, welchen OIDC/OAuth2-Anbieter Sie integrieren möchten. Hier ist eine Liste der gängigen Anbieter - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Bitte senden Sie uns die Informationen, die für die Konfiguration erforderlich sind (in der Regel sind es Anbietername, issuer, client_id, client_secret, aber das hängt vom Anbieter ab)

Sygnal: Sind Sie sicher, dass Sie es wollen? Es ist ein Push Gateway, das nur für Matrix-Client-App-Entwickler nutzbar ist. Sie können es also nicht verwenden, wenn Sie Ihre mobile Matrix-App nicht entwickeln. Wenn Sie es hinzufügen möchten, geben Sie bitte die folgenden Informationen an: App ID(-s) (z.B. org.matrix.app), FCM api key, und/oder APNS certificate (falls verwendet)

BorgBackup: geben Sie bitte die gewünschte Repository-URL an (user@host:repo). Wir werden einen SSH-Schlüssel und eine Verschlüsselungspassphrase auf Ihrem Server generieren. Wir senden Ihnen den öffentlichen Teil des generierten SSH-Schlüssels zu. Sie müssen diesen SSH-Schlüssel zu Ihrem Provider hinzufügen.

Jitsi: Sind Sie sicher, dass Sie es wollen? Sie erhalten die jitsi-Integration standardmäßig mit einer öffentlichen Instanz. Das von uns angebotene jitsi ist eine selbst gehostete Version. Denken Sie daran, dass Jitsi die Anforderungen an die Rechenleistung deutlich erhöht.

Server: Bitte erstellen Sie einen x86/amd64 VPS mit einer beliebigen Debian-basierten Distribution. Minimale komfortable Konfiguration für einen einfachen Matrix-Server: 1vCPU, 2GB RAM.
 Fügen Sie unsere SSH-Schlüssel (https://etke.cc/ssh.key) zu Ihrem Server hinzu, senden Sie uns die IP-Adresse Ihres Servers, den Benutzernamen (mit der Berechtigung, sudo aufzurufen) und das Passwort (falls festgelegt).

```


___

```yaml

DNS - bitte fügen Sie die folgenden Einträge hinzu:

@	A record	$SERVER_IP4
matrix	A record	$SERVER_IP4
buscarron	CNAME record	matrix.example.com
cinny	CNAME record	matrix.example.com
dimension	CNAME record	matrix.example.com
element	CNAME record	matrix.example.com
etherpad	CNAME record	matrix.example.com
hydrogen	CNAME record	matrix.example.com
jitsi	CNAME record	matrix.example.com
ntfy	CNAME record	matrix.example.com
stats	CNAME record	matrix.example.com
sygnal	CNAME record	matrix.example.com
matrix	MX record	matrix.example.com
matrix	TXT record	v=spf1 ip4:$SERVER_IP4 -all
_dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
```

