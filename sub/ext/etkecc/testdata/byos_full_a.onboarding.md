# links

* homeserver: https://matrix.example.com
* synapse-admin: https://matrix.example.com/synapse-admin
* etherpad admin: https://dimension.example.com/etherpad/admin
* ssh git: ssh matrix.example.com:23231
* buscarron: https://buscarron.example.com
* cinny: https://cinny.example.com
* dimension: https://dimension.example.com
* element-web: https://element.example.com
* go-neb: https://goneb.example.com
* hydrogen: https://hydrogen.example.com
* jitsi: https://jitsi.example.com
* kuma: https://kuma.example.com
* languagetool: https://languagetool.example.com
* miniflux: https://miniflux.example.com
* miounne: https://miounne.example.com
* ntfy: https://ntfy.example.com
* radicale: https://radicale.example.com
* stats: https://stats.example.com
* sygnal: https://sygnal.example.com

# bots

* buscarron: @buscarron:example.com
* honoroit: @honoroit:example.com
* miounne: @miounne:example.com
* mjolnir: @mjolnir:example.com
* postmoogle: @postmoogle:example.com
* reminder-bot: @reminder:example.com

# bridges

* discord: @discordbot:example.com
* facebook: @facebookbot:example.com
* googlechat: @googlechatbot:example.com
* groupme: @_groupmepuppet_bot:example.com
* instagram: @instagrambot:example.com
* irc: @heisenbridge:example.com
* linkedin: @linkedinbot:example.com
* signal: @signalbot:example.com
* skype: @skypebridgebot:example.com
* slack: @_slackpuppet_bot:example.com
* steam: @_steampuppet_bot:example.com
* telegram: @telegrambot:example.com
* twitter: @twitterbot:example.com
* webhooks: @_webhook:example.com
* whatsapp: @whatsappbot:example.com

> https://etke.cc/de/bridges - Authentifizierungsanweisungen

# Anmeldeinformationen

* mxid: @test:example.com
* username: test
* password: TODO
* etherpad admin password: TODO
* grafana password: TODO
* matrix-corporal api password: TODO
* radicale password: TODO

> im Falle irgendwelcher Probleme: https://etke.cc/de/help

# Bezahlung

Bitte, [kaufen Sie das Einrichtungselement](https://etke.cc/setup) und Beitritt zur **Maintenance+Email** Mitgliedschaft auf https://etke.cc/membership.

> **HINWEIS**: Alle Preise basieren auf dem [Pay What You Want]-Modell (https://de.wikipedia.org/wiki/Pay_what_you_want).

# Schritte nach dem Einrichten

### email2matrix

1. Erstellen Sie einen Nicht-Administrator-Matrix-Benutzer mit dem Benutzernamen `email2matrix` und einem sicheren Passwort
2. Erstellen Sie unverschlüsselte Matrixräume für die gewünschten Postfächer (1 Raum = 1 Postfach) und laden Sie den email2matrix-Benutzer zu allen Räumen ein.
3. Melden Sie sich mit dem email2matrix-Benutzer an und akzeptieren Sie die Einladungen, die Sie verschickt haben (ja, das ist ein sehr manueller Prozess - deshalb haben wir Sie gefragt, ob Sie das wirklich wollen)
4. Senden Sie das email2matrix-Passwort und Paare von (Mailboxname + Matrix-Raum-ID) an den @support:etke.cc (eg: info@matrix.example.com = !gqlCuoCdhufltluRXk:example.com)

### postmoogle

1. Befehl `!pm dkim` in einem beliebigen Raum mit postmoogle ausführen
2. Fügen Sie einen neuen TXT-DNS-Eintrag mit Key/From/Subdomain = `postmoogle._domainkey.matrix` und Value/To = Signatur aus der Ausgabe von `!pm dkim` hinzu

### etherpad

1. Öffnen Sie den Integrationsmanager in einer beliebigen Element-Client-Anwendung
2. Klicken Sie auf das Zahnradsymbol und gehen Sie auf die Registerkarte Widgets
3. Öffnen Sie die Etherpad-Widget-Konfiguration und ersetzen Sie `scalar.vector.im` durch `dimension.example.com`

### buscarron

1. Erstellen Sie (verschlüsselte) Matrixräume für die gewünschten Formulare und laden Sie den buscarron-Benutzer zu allen Räumen ein
2. Senden Sie die Liste (Formularname, Matrix-Raum-ID, Umleitungs-URL nach Übermittlung) an @support:etke.cc

### honoroit

1. Erstellen Sie einen Matrixraum (Verschlüsselung unterstützt) und laden Sie den honoroit-Benutzer dazu ein
2. Senden Sie die Raum-ID an @support:etke.cc

Wenn Sie die Nachrichten von honoroit ändern möchten, senden Sie die gewünschten Texte an @support:etke.cc (Sie können https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go als Referenz verwenden)