# ссылки

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

# боты

* buscarron: @buscarron:higenjitsuteki.etke.host
* honoroit: @honoroit:higenjitsuteki.etke.host
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
* password: TODO64
* etherpad admin password: TODO64
* grafana password: TODO64

> в случае проблем: https://etke.cc/ru/help

# после установки

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