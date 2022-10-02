# ссылки

* homeserver: https://matrix.higenjitsuteki.etke.host
* synapse-admin: https://matrix.higenjitsuteki.etke.host/synapse-admin
* etherpad admin: https://dimension.higenjitsuteki.etke.host/etherpad/admin
* ssh git: ssh matrix.higenjitsuteki.etke.host:23231
* buscarron: https://buscarron.higenjitsuteki.etke.host
* cinny: https://cinny.higenjitsuteki.etke.host
* dimension: https://dimension.higenjitsuteki.etke.host
* element-web: https://element.higenjitsuteki.etke.host
* go-neb: https://goneb.higenjitsuteki.etke.host
* hydrogen: https://hydrogen.higenjitsuteki.etke.host
* jitsi: https://jitsi.higenjitsuteki.etke.host
* kuma: https://kuma.higenjitsuteki.etke.host
* languagetool: https://languagetool.higenjitsuteki.etke.host
* miniflux: https://miniflux.higenjitsuteki.etke.host
* miounne: https://miounne.higenjitsuteki.etke.host
* ntfy: https://ntfy.higenjitsuteki.etke.host
* radicale: https://radicale.higenjitsuteki.etke.host
* stats: https://stats.higenjitsuteki.etke.host
* sygnal: https://sygnal.higenjitsuteki.etke.host

# боты

* buscarron: @buscarron:higenjitsuteki.etke.host
* honoroit: @honoroit:higenjitsuteki.etke.host
* miounne: @miounne:higenjitsuteki.etke.host
* mjolnir: @mjolnir:higenjitsuteki.etke.host
* reminder-bot: @reminder:higenjitsuteki.etke.host

# интеграции

* discord: @discordbot:higenjitsuteki.etke.host
* facebook: @facebookbot:higenjitsuteki.etke.host
* googlechat: @googlechatbot:higenjitsuteki.etke.host
* groupme: @_groupmepuppet_bot:higenjitsuteki.etke.host
* instagram: @instagrambot:higenjitsuteki.etke.host
* irc: @heisenbridge:higenjitsuteki.etke.host
* linkedin: @linkedinbot:higenjitsuteki.etke.host
* signal: @signalbot:higenjitsuteki.etke.host
* skype: @skypebridgebot:higenjitsuteki.etke.host
* slack: @_slackpuppet_bot:higenjitsuteki.etke.host
* steam: @_steampuppet_bot:higenjitsuteki.etke.host
* telegram: @telegrambot:higenjitsuteki.etke.host
* twitter: @twitterbot:higenjitsuteki.etke.host
* webhooks: @_webhook:higenjitsuteki.etke.host
* whatsapp: @whatsappbot:higenjitsuteki.etke.host

> https://etke.cc/ru/bridges - инструкции по авторизации

# доступы

* mxid: @test:higenjitsuteki.etke.host
* username: test
* password: TODO
* etherpad admin password: TODO
* grafana password: TODO
* matrix-corporal api password: TODO
* radicale password: TODO

> в случае проблем: https://etke.cc/ru/help

# после установки

### email2matrix

1. Создайте обычного пользователя matrix с логином `email2matrix` и надежным паролем
2. На каждый почтовый ящик, который Вы хотите добавить, создайте незашифрованную matrix комнату и пригласите пользователя email2matrix в каждую из них
3. Войдите пользователем email2matrix и примите все приглашения, которые были отправлены ранее (да, вот потому мы и спрашивали, точно ли Вы хотите это)
4. отправьте пароль пользователя email2matrix и пары почтовый ящик - ID matrix комнаты на @support:etke.cc (eg: info@matrix.higenjitsuteki.etke.host = !gqlCuoCdhufltluRXk:higenjitsuteki.etke.host)

### etherpad

1. Откройте менеджер интеграций в любом приложении Element
2. Нажмите на иконку шестеренки и перейдите в таб Widgets
3. Откройте конфигурацию виджета Etherpad и замените `scalar.vector.im` на `dimension.higenjitsuteki.etke.host`

### buscarron

1. Создайте (зашифрованные) комнаты matrix для всех форм, которые Вы собираетесь сделать и пригласите пользователя buscarron в каждую из них
2. Отправьте список (название формы, id комнаты matrix, URL для перенаправления пользователя после отправки) на @support:etke.cc

### honoroit

1. Создайте matrix комнату (шифрование поддерживается) и пригласите пользователя honoroit в нее
2. Отправьте ID этой комнаты на @support:etke.cc

Если Вы хотите изменить текст сообщений, отправьте желаемый вариант на @support:etke.cc (используйте https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go в качесте справки)