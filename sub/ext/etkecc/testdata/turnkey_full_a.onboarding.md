# links

* homeserver: https://matrix.higenjitsuteki.etke.host
* etherpad admin: https://dimension.higenjitsuteki.etke.host/etherpad/admin
* buscarron: https://buscarron.higenjitsuteki.etke.host
* cinny: https://cinny.higenjitsuteki.etke.host
* dimension: https://dimension.higenjitsuteki.etke.host
* element-web: https://element.higenjitsuteki.etke.host
* hydrogen: https://hydrogen.higenjitsuteki.etke.host
* jitsi: https://jitsi.higenjitsuteki.etke.host
* ntfy: https://ntfy.higenjitsuteki.etke.host
* stats: https://stats.higenjitsuteki.etke.host
* sygnal: https://sygnal.higenjitsuteki.etke.host

# bots

* buscarron: @buscarron:higenjitsuteki.etke.host
* honoroit: @honoroit:higenjitsuteki.etke.host
* reminder-bot: @reminder:higenjitsuteki.etke.host

# bridges

* discord: @discordbot:higenjitsuteki.etke.host
* facebook: @facebookbot:higenjitsuteki.etke.host
* googlechat: @googlechatbot:higenjitsuteki.etke.host
* groupme: @_groupmepuppet_bot:higenjitsuteki.etke.host
* instagram: @instagrambot:higenjitsuteki.etke.host
* irc: @heisenbridge:higenjitsuteki.etke.host
* linkedin: @linkedinbot:higenjitsuteki.etke.host
* signal: @signalbot:higenjitsuteki.etke.host
* skype: @skypebridgebot:higenjitsuteki.etke.host
* slack: @slackbot:higenjitsuteki.etke.host
* steam: @_steampuppet_bot:higenjitsuteki.etke.host
* telegram: @telegrambot:higenjitsuteki.etke.host
* twitter: @twitterbot:higenjitsuteki.etke.host
* webhooks: @hookshot:higenjitsuteki.etke.host
* whatsapp: @whatsappbot:higenjitsuteki.etke.host

> https://etke.cc/en/help/bridges - auth instructions

# credentials

* mxid: @test:higenjitsuteki.etke.host
* username: test
* password: TODO64
* etherpad admin password: TODO64
* grafana password: TODO64

> in case of any issues: https://etke.cc/en/help

# steps after the setup

### etherpad

1. Open integration manager in any element client app
2. Click on the sprocket icon and go to the Widgets tab
3. Open the Etherpad widget configuration and replace `scalar.vector.im` with `dimension.higenjitsuteki.etke.host`

### buscarron

1. Create (encrypted) matrix room(-s) for the forms you want to have and invite the buscarron user to all of them
2. Send the list of (form name, matrix room id, redirect URL after submission) to the @support:etke.cc

### honoroit

1. Create a matrix room (encryption supported) and invite the honoroit user into it
2. Send the room id to the @support:etke.cc

if you want to change honoroit's messages, send the texts you want to use to @support:etke.cc (you can use https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go as a reference)