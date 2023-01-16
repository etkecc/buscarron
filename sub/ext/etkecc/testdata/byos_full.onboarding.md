# links

* homeserver: https://matrix.example.com
* etherpad admin: https://dimension.example.com/etherpad/admin
* buscarron: https://buscarron.example.com
* cinny: https://cinny.example.com
* dimension: https://dimension.example.com
* element-web: https://element.example.com
* hydrogen: https://hydrogen.example.com
* jitsi: https://jitsi.example.com
* ntfy: https://ntfy.example.com
* stats: https://stats.example.com
* sygnal: https://sygnal.example.com

# bots

* buscarron: @buscarron:example.com
* honoroit: @honoroit:example.com
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
* slack: @slackbot:example.com
* steam: @_steampuppet_bot:example.com
* telegram: @telegrambot:example.com
* twitter: @twitterbot:example.com
* webhooks: @hookshot:example.com
* whatsapp: @whatsappbot:example.com

> https://etke.cc/de/help/bridges - Authentifizierungsanweisungen

# Anmeldeinformationen

* mxid: @test:example.com
* username: test
* password: TODO64
* etherpad admin password: TODO64
* grafana password: TODO64

> im Falle irgendwelcher Probleme: https://etke.cc/de/help

# Schritte nach dem Einrichten

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