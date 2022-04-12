# links

* homeserver: https://matrix.example.com
* synapse-admin: https://matrix.example.com/synapse-admin
* etherpad admin: https://dimension.example.com/etherpad/admin
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
* radicale: https://radicale.example.com
* stats: https://stats.example.com
* sygnal: https://sygnal.example.com


# bots

* honoroit: @honoroit:example.com
* miounne: @miounne:example.com
* mjolnir: @mjolnir:example.com
* reminder-bot: @reminder:example.com


# bridges

* discord: @_discordpuppet_bot:example.com
* facebook: @facebookbot:example.com
* googlechat: @googlechatbot:example.com
* groupme: @_groupmepuppet_bot:example.com
* instagram: @instagrambot:example.com
* irc: @heisenbridge:example.com
* linkedin: @linkedinbot:example.com
* signal: @signalbot:example.com
* slack: @_slackpuppet_bot:example.com
* steam: @_steampuppet_bot:example.com
* telegram: @telegrambot:example.com
* twitter: @twitterbot:example.com
* webhooks: @_webhook:example.com
* whatsapp: @whatsappbot:example.com


> https://etke.cc/bridges - auth instructions

# credentials

* mxid: @test:example.com
* username: test
* password: TODO
* etherpad admin password: TODO
* grafana password: TODO
* matrix-corporal api password: TODO
* radicale password: TODO


> in case of any issues: @support:etke.cc

# payment

Please, [buy the Setup item](https://etke.cc/setup) and join the **Maintenance+Email** membership on [https://etke.cc/membership](https://etke.cc/membership).

> **NOTE**: all prices are based on [Pay What You Want](https://en.wikipedia.org/wiki/Pay_what_you_want) model.

# steps after the setup

### email2matrix

1. Add new MX record on matrix subdomain that will point... to matrix.example.com (looks odd, but some mail servers will not send email to your server without it)
2. Create a non-admin matrix user with username `email2matrix` and secure password
3. Create un-encrypted matrix rooms for mailboxes you want to have (1 room = 1 mailbox) and invite email2matrix user in all of them
4. Login with email2matrix user and accept the invites you sent (yes, that's why I asked you if you really want it)
5. Send email2matrix password and pairs of mailbox name - matrix room id to the @support:etke.cc (eg: info@matrix.example.com = !gqlCuoCdhufltluRXk:example.com)

### etherpad

1. Open integration manager in any element client app
2. Click on the sprocket icon and go to the Widgets tab
3. Open the Etherpad widget configuration and replace `scalar.vector.im` with `dimension.example.com`

### honoroit

1. Create a matrix room (encryption supported) and invite the honoroit user into it
2. Send the room id to the @support:etke.cc

if you want to change honoroit messages, send the texts you want to set to the @support:etke.cc (you can use https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go as reference)

