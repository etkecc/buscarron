# ссылки

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
* radicale: https://radicale.example.com
* stats: https://stats.example.com
* sygnal: https://sygnal.example.com


# боты

* honoroit: @honoroit:example.com
* miounne: @miounne:example.com
* mjolnir: @mjolnir:example.com
* reminder-bot: @reminder:example.com


# интеграции

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


> https://etke.cc/bridges - инструкции по авторизации

# доступы

* mxid: @test:example.com
* username: test
* password: TODO
* etherpad admin password: TODO
* grafana password: TODO
* matrix-corporal api password: TODO
* radicale password: TODO


> в случае проблем: @support:etke.cc

# после установки

### email2matrix

1. Добавь новую MX запись на поддомене matrix, которая будет указывать... на matrix.example.com (выглядит странно, но некторые почтовые сервера не будут доставлять письма на твой сервер без нее)
2. Создай обычного пользователя matrix с логином `email2matrix` и надежным паролем
3. На каждый почтовый ящик, который ты хочешь добавить, создай незашифрованную matrix комнату и пригласи пользователя email2matrix в каждую из них
4. Войди пользователем email2matrix и прими все приглашения, которые были отправлены тобой ранее (да, вот потому мы и спрашивали, точно ли ты хочешь это)
5. Отправь пароль пользователя email2matrix и пары почтовый ящик - ID matrix комнаты на @support:etke.cc (eg: info@matrix.example.com = !gqlCuoCdhufltluRXk:example.com)

### etherpad

1. Открой менеджер интеграций в любом приложении Element
2. Нажми на иконку шестеренки и перейди в таб Widgets
3. Открой конфигурацию виджета Etherpad и замени `scalar.vector.im` на `dimension.example.com`

### honoroit

1. Создай matrix комнату (шифрование поддерживается) и пригласи пользователя honoroit в нее
2. Отправь ID этой комнаты на @support:etke.cc

Если ты хочешь изменить текст сообщений, отправь желаемый вариант на @support:etke.cc (используй https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go в качесте справки)

