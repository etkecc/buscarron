```yaml
Hallo zusammen,
Wir haben Ihre Bestellung erhalten und haben einige Fragen vor der Einrichtung.

Wir sehen, dass Sie etwas auf Ihrer Basisdomain haben. In diesem Fall sollten Sie die folgenden HTTPS-Weiterleitungen (HTTP 301) hinzufügen:
* https://example.com/.well-known/matrix/server -> https://matrix.example.com/.well-known/matrix/server
* https://example.com/.well-known/matrix/client -> https://matrix.example.com/.well-known/matrix/client

Reminder bot: Was ist Ihre Zeitzone (IANA)? Wie America/Chicago, Asia/Seoul, oder Europe/Berlin

Honoroit: Sind Sie sicher, dass Sie es wollen? Es ist ein Helpdesk-Bot mit Ende-Zu-Ende-Verschlüsselungsunterstützung. Bitte prüfen Sie https://gitlab.com/etke.cc/honoroit, bevor Sie sich entscheiden.

Telegram: gehen Sie bitte auf https://my.telegram.org/apps und erstellen Sie eine neue App. Teilen Sie mir die API-ID und den Hash mit

SMTP relay: Bitte wählen Sie einen geeigneten E-Mail-Anbieter (große Anbieter wie Gmail oder Outlook sperren Sie für automatisierte E-Mails, daher müssen Sie einen Dienst finden, der den Versand von Verifizierungs-E-Mails erlaubt. Optional bieten wir einen solchen Dienst an). Bitte senden Sie uns einen SMTP Host, einen SMTP STARTTLS Port, ein SMTP Login, ein SMTP Passwort und eine SMTP E-Mail (normalerweise sind Login und E-Mail dasselbe, aber das hängt vom Anbieter ab).

Prometheus+Grafana: Sind Sie sicher, dass Sie das wollen? Cloud-Anbieter bieten in der Regel ein Dashboard mit Server-Statistiken, warum also nicht stattdessen dieses Dashboard verwenden? Ein Prometheus+Grafana-Stack bietet einige interne Matrix-Statistiken (z. B. die Anzahl der Ereignisse), ist aber zu viel des Guten, wenn Sie nur die Serverauslastung sehen wollen.

Uptime Kuma: Sind Sie sicher, dass Sie es wollen? Es handelt sich um ein selbstgehostetes Überwachungsprogramm, das nicht in Matrix integriert ist, wie z. B. 'opensource UptimeRobot'. Bitte prüfen Sie https://github.com/louislam/uptime-kuma, bevor Sie sich entscheiden.

Radicale: Sind Sie sicher, dass Sie es wollen? Es handelt sich um einen CalDAV/CardDAV-Server, der nicht in Matrix integriert ist. Bitte prüfen Sie https://radicale.org/, bevor Sie sich entscheiden.

Miniflux: Sind Sie sicher, dass Sie es wollen? Es ist ein RSS-Reader, der nicht in Matrix integriert ist. Bitte prüfen Sie https://miniflux.app, bevor Sie sich entscheiden.

Languagetool: Sind Sie sicher, dass Sie es wollen? Es handelt sich um einen 'Open-Source-Grammarly'-Server, der ~30 GB Speicherplatz für n-Gramme benötigt und nicht in Matrix integriert ist. Bitte prüfen Sie https://languagetool.org, bevor Sie sich entscheiden.

Soft-Serve: Sind Sie sicher, dass Sie es wollen? Es ist ein Git-Hosting das über SSH zugänglich ist, das nicht in Matrix integriert ist. Bitte prüfen Sie https://github.com/charmbracelet/soft-serve, bevor Sie sich entscheiden.

WireGuard und dnsmasq: sind Sie sicher, dass Sie sie wollen? WireGuard ist ein VPN (nicht in Matrix integriert) und dnsmasq ist ein lokaler DNS-Server, der mit Adblock-Listen (wie pi-hole) ausgestattet ist, die ziemlich gut in WireGuard integriert sind. Bitte prüfen Sie https://wireguard.com, bevor Sie sich entscheiden. Wenn Sie es immer noch wollen, senden Sie uns bitte eine Liste von Bezeichnungen, die Sie den generierten Client-Schlüsseln zuweisen wollen (nur um Dateinamen festzulegen, also auch '1,2,3...' ist OK)

Etherpad (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie das wollen? Es handelt sich um einen selbst gehosteten kollaborativen Editor, der nur in Verbindung mit Dimension verwendet werden kann (auch dazu wurde eine Frage gestellt). Denken Sie daran, dass Sie mit dem Standard-Integrationsmanager ohnehin ein Etherpad erhalten, das jedoch von Element Inc. gehostet wird (Entwickler der Element-Clientanwendungen).

Dimension (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie das wollen? Es ist ein selbst gehosteter Integrationsmanager. Sie erhalten standardmäßig einen Integrationsmanager mit jeder Element-Client-Anwendung. Bitte prüfen Sie https://github.com/turt2live/matrix-dimension bevor Sie sich entscheiden.

Website (nur mit Abonnement oder Komplettpaket): Um eine statische Website einzurichten, müssen Sie Ihre Basisdomäne (den @ DNS eintrag) auf die IP-Adresse des Matrix-Servers verweisen, und der Quellcode der Website muss in einem öffentlichen git repo verfügbar sein. Sind Sie sicher, dass Sie das wollen? Wenn ja, geben Sie bitte die URL des Website-Repositorys, den Befehl (-s) zum Erstellen der Website und den Ordner an, in dem die Build-Distribution gespeichert ist (normalerweise public oder dist).

buscarron (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Es ist ein Bot, der Webformulare (HTML/HTTP POST) empfängt und sie an (verschlüsselte) Matrixräume sendet. Bitte prüfen Sie https://gitlab.com/etke.cc/buscarron, bevor Sie sich entscheiden.

SSO (nur mit Abonnement oder Komplettpaket): Sie haben nicht erwähnt, welchen OIDC/OAuth2-Anbieter Sie integrieren möchten. Hier ist eine Liste der gängigen Anbieter - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Bitte senden Sie uns die Informationen, die für die Konfiguration erforderlich sind (in der Regel sind es Anbietername, issuer, client_id, client_secret, aber das hängt vom Anbieter ab)

Sygnal (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Es ist ein Push Gateway, das nur für Matrix-Client-App-Entwickler nutzbar ist. Sie können es also nicht verwenden, wenn Sie Ihre mobile Matrix-App nicht entwickeln. Wenn Sie es hinzufügen möchten, geben Sie bitte die folgenden Informationen an: App ID(-s) (z.B. org.matrix.app), FCM api key, und/oder APNS certificate (falls verwendet)

BorgBackup (nur mit Abonnement oder Komplettpaket): geben Sie bitte die gewünschte Repository-URL an (user@host:repo). Wir werden einen SSH-Schlüssel und eine Verschlüsselungspassphrase auf Ihrem Server generieren. Wir senden Ihnen den öffentlichen Teil des generierten SSH-Schlüssels zu. Sie müssen diesen SSH-Schlüssel zu Ihrem Provider hinzufügen.

email2matrix (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Es handelt sich um einen Einweg-SMTP-Server für den Empfang von E-Mails in einem Matrixraum. Es ist ziemlich schwierig, ihn selbst einzurichten, da es keinen einfachen Weg gibt, ihn zu konfigurieren. Wir müssen mit Ihnen zusammenarbeiten, um die Konfiguration sowohl als Matrix-Homeserver-Benutzer (von Ihnen durchgeführt, da wir keine Benutzer auf Ihrem Homeserver haben und keinen Zugriff auf Ihre Daten innerhalb von Matrix haben) als auch für das System (von uns durchgeführt, da die Konfiguration dieses Tools in Konfigurationsdateien auf der VM/VPS-Festplatte gespeichert wird) durchzuführen.

Jitsi (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Sie erhalten die jitsi-Integration standardmäßig mit einer öffentlichen Instanz. Das von uns angebotene jitsi ist eine selbst gehostete Version. Denken Sie daran, dass Jitsi die Anforderungen an die Rechenleistung deutlich erhöht.

ma1sd (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Es handelt sich um veraltete Software, die früher als Stub - Identitätsserver - verwendet wurde, seit einiger Zeit nicht mehr gewartet wird und in den meisten Fällen keine Vorteile bietet (es sei denn, Sie möchten LDAP-Authentifizierung oder Twilio-Telefonnummernüberprüfung hinzufügen).

matrix-registration (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Es handelt sich um veraltete Software - ein Workaround, der verwendet wurde, um Matrix eine einladungsbasierte Registrierungsunterstützung hinzuzufügen, weil das Protokoll dies nicht unterstützte, aber jetzt können Sie stattdessen die eingebaute Einladungs-Token-Funktion verwenden: https://matrix-org.github.io/synapse/latest/usage/administration/admin_api/registration_tokens.html

Miounne (nur mit Abonnement oder Komplettpaket): Sind Sie sicher, dass Sie es wollen? Miounne ist eine veraltete Software - eine Brücke zwischen externen Diensten (wie HTML/HTTP-Formulare, Matrix-Registrierung, buymeacoffee, etc.) und Matrix. Bitte prüfen Sie https://gitlab.com/etke.cc/miounne, bevor Sie sich entscheiden. Wenn Sie es immer noch wollen, schicken Sie mir bitte eine Konfiguration, die ich anwenden kann (nein, es gibt keine 'Standardkonfiguration'. Nein, es gibt keine 'gute Konfiguration'. Nein, wir stellen keine Konfigurationsvorlagen zur Verfügung. Es liegt ganz bei Ihnen)

Server: Bitte erstellen Sie einen VPS mit einer beliebigen Debian-basierten Distribution. Minimale komfortable Konfiguration für einen einfachen Matrix-Server: 1vCPU, 2GB RAM.
 Fügen Sie unsere SSH-Schlüssel (https://etke.cc/ssh.key) zu Ihrem Server hinzu, senden Sie uns die IP-Adresse Ihres Servers, den Benutzernamen (mit der Berechtigung, sudo aufzurufen) und das Passwort (falls festgelegt).

```


___

```yaml

DNS - bitte fügen Sie die folgenden Einträge hinzu:
matrix	A record	server IP
buscarron	CNAME record	matrix.example.com
cinny	CNAME record	matrix.example.com
dimension	CNAME record	matrix.example.com
element	CNAME record	matrix.example.com
goneb	CNAME record	matrix.example.com
hydrogen	CNAME record	matrix.example.com
jitsi	CNAME record	matrix.example.com
kuma	CNAME record	matrix.example.com
languagetool	CNAME record	matrix.example.com
miniflux	CNAME record	matrix.example.com
miounne	CNAME record	matrix.example.com
ntfy	CNAME record	matrix.example.com
radicale	CNAME record	matrix.example.com
stats	CNAME record	matrix.example.com
sygnal	CNAME record	matrix.example.com
matrix	MX record	matrix.example.com
matrix	TXT record	v=spf1 ip4:SERVER_IP -all
_dmarc.matrix	TXT record	v=DMARC1; p=quarantine;
```

