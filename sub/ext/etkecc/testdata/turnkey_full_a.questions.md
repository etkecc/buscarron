```yaml
Привет,
Мы получили Ваш заказ и хотели задать несколько вопросов перед установкой.

Мы видим, что у Вас что-то настроено на основном домене. В этом случае Вам нужно будет добавить HTTPS редиректы (HTTP 301):
* https://higenjitsuteki.etke.host/.well-known/matrix/server -> https://matrix.higenjitsuteki.etke.host/.well-known/matrix/server
* https://higenjitsuteki.etke.host/.well-known/matrix/client -> https://matrix.higenjitsuteki.etke.host/.well-known/matrix/client

Reminder bot: Какой у Вас часовой пояс (в формате IANA)? Например, America/Chicago, Asia/Seoul или Europe/Berlin

Honoroit: Вы точно хотите это? Это хелпдеск бот с поддержкой e2e шифрования. Пожалуйста, посмотрите на https://gitlab.com/etke.cc/honoroit и решите.

Telegram: пожалуйста, сходите на https://my.telegram.org/apps и создайте новое приложение. Отправьте нам API ID и Hash

Uptime Kuma: Вы точно хотите это? Простой мониторинг сайтов на Вашем сервере, без интеграции с matrix, что-то вроде 'open source UptimeRobot'. Пожалуйста, откройте https://github.com/louislam/uptime-kuma и решите.

Radicale: Вы точно хотите это? CalDav/CardDav сервер, без интеграции с matrix. Пожалуйста, откройте https://radicale.org и решите.

Miniflux: Вы точно хотите это? Читалка RSS, без интеграции с matrix. Пожалуйста, откройте https://miniflux.app и решите.

Languagetool: Вы точно хотите это? Что-то вроде сервера 'open source Grammarly', требует около 30ГБ места на диске для н-грамм, без интеграции с matrix. Пожалуйста, откройте https://languagetool.org и решите.

Soft-Serve: Вы точно хотите это? Git хостинг, достуный только по ssh, без интеграции с matrix. Пожалуйста, откройте https://github.com/charmbracelet/soft-serve и решите.

WireGuard и dnsmasq: Вы точно хотите их? WireGuard это VPN (без интеграции с matrix), а dnsmasq - локальный DNS сервер с блокировщиком рекламы (как pi-hole), который хорошо интегрирован с WireGuard. Пожалуйста, откройте https://wireguard.com и решите. Если Вы хотите добавить их, пожалуйста, отправьте нам список названий, которые Вы хотите использовать для клиентов (просто для имен файлов, так что даже '1,2,3,...' подойдет)

Etherpad (только с подпиской или хостингом): Вы точно хотите это? Это редактор текста для совместной работы на Вашем сервере, работает только в связке с dimension (добавили вопрос и о нем). Etherpad и так доступен Вам по умолчанию в менеджере интеграций в любом приложении Element (хостится на серверах Element Inc., разработчики приложений Element).

Dimension (только с подпиской или хостингом): Вы точно хотите это? Это менеджер интеграций на Вашем сервере. По умолчанию у Вас и так будет доступен стандартный менеджер интеграций. Пожалуйста, посмотрите на https://github.com/turt2live/matrix-dimension и решите

Website (только с подпиской или хостингом): чтобы задеплоить статичный вебсайт, Вам придется настроить свой базовый домен на сервер matrix (@ DNS запись), а исходники самого вебсайта должны быть доступны в публичном git репозитории. Вы точно хотите этого? Если да, пожалуйста, отправьте нам адрес git репозитория Вашего статичного сайта, список команд для сборки и в какой директории будет собран артефакт (обычно это public или dist).

buscarron (только с подпиской или хостингом): Вы точно хотите это? Это бот, который принимает отправку веб форм (HTML/HTTP POST) и отправляет их в (зашифрованные) Matrix комнаты. Пожалуйста, посмотрите на https://gitlab.com/etke.cc/buscarron и решите.

SSO (только с подпиской или хостингом): Мы не получили информацию о том, какого OIDC/OAuth2 провайдера Вы хотите интегрировать, так что вот список популярных провайдеров - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Пожалуйста, отправьте нам требуемую информацию для настройки (обычно это название провайдера, issuer, client_id и client_secret, но все зависит от выбранного провайдера)

Sygnal (только с подпиской или хостингом): Вы точно хотите это? Это push gateway, нужен только разработчикам мобильных приложений matrix, так что он будет бесполезен, если Вы не разрабатываете собственное мобильное приложение для matrix. Если Вы все же хотите его добавить, пожалуйста, отправьте следующую информацию: ID приложений (например org.matrix.app), API ключ FCM и/или сертификат APNS (если используется)

BorgBackup (только с подпиской или хостингом): пожалуйста, предоставьте желаемый адрес репозитория (user@host:repo). Мы сгенерируем ssh ключ и пароль шифрования на стороне Вашего сервера и отправим Вам публичную часть ключа. Вам нужно будет добавить этот ключ на стороне Вашего borg провайдера.

email2matrix (только с подпиской или хостингом): Вы точно хотите это? Это односторонний SMTP сервер для получения имейлов в комнате matrix и настроить его Вам будет крайне непросто, нам нужно будет скооперироваться с тобой, чтобы подготовить конфигурацию с обеих сторон - matrix (Вы - у нас нет matrix ползователей на Вашем сервере, потому и нет доступа к данным внутри matrix) и системы (мы - потому как все настройки email2matrix надо править прямо в файлах конфигурации на диске).

Jitsi (только с подпиской или хостингом): Вы точно хотите это? По умолчанию у Вас будет интегрирован публичный сервер jitsi, а мы предлагаем то же самое установить на Ваш сервер. Имейте в виду, что jitsi на Вашем сервере сильно увеличивает требования к конфигурации сервера.

ma1sd (только с подпиской или хостингом): Вы точно хотите это? Это устаревшее решение, изначально использовалось как заглушка - identity сервер, не обновлявшийся уже долгое время и не имеющий никаких плюсов в большинстве случаев (за исключением ситуаций, когда Вам нужна авторизация LDAP или подтверждение телефона через Twilio).

matrix-registration (только с подпиской или хостингом): Вы точно хотите это? Это устаревшее решение - костыль, чтобы добавить регистрацию по приглашениям в matrix, потому что раньше протокол не поддерживал их, но теперь по умолчанию включено встроенное решение: https://matrix-org.github.io/synapse/latest/usage/administration/admin_api/registration_tokens.html

Miounne (только с подпиской или хостингом): Вы точно хотите это? Это устаревшее решение - интеграция между внешними сервисами (например, HTML/HTTP формы, matrix-registration, buymeacoffee и т.д.) и matrix. Пожалуйста, откройте https://gitlab.com/etke.cc/miounne и решите. Если Вы все же хотите добавить miounne, пожалуйста, отправьте нам желаемую конфигурацию (нет, 'конфигурации по умолчанию' не существует. Нет, 'хорошей конфигурации' не существует. Нет, мы не предоставляем шаблоны конфигурации, это все на Вас)

q_turnkey_ssh

```


___

```yaml
export HETZNER_SERVER_IP=SERVER_IP
curl -X "POST" "https://dns.hetzner.com/api/v1/records/bulk" -H "Content-Type: application/json" -H "Auth-API-Token: $HETZNER_API_TOKEN" -d "{\"records\":[{\"name\":\"higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"A\",\"value\":\"$HETZNER_SERVER_IP\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"10 aspmx1.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"higenjitsuteki\",\"type\":\"MX\",\"value\":\"20 aspmx2.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"autoconfig.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"autoconfig.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"_autodiscover._tcp.higenjitsuteki\",\"type\":\"SRV\",\"value\":\"0 1 443 autodiscover.migadu.com\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"key1._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key1.higenjitsuteki._domainkey.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"key2._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key2.higenjitsuteki._domainkey.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"key3._domainkey.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"key3.higenjitsuteki._domainkey.migadu.com.\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"_dmarc.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 include:spf.migadu.com -all\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"buscarron.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"cinny.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"dimension.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"element.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"goneb.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"hydrogen.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"jitsi.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"kuma.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"languagetool.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"miniflux.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"miounne.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"ntfy.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"radicale.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"stats.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"sygnal.higenjitsuteki\",\"type\":\"CNAME\",\"value\":\"matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"MX\",\"value\":\"0 matrix.higenjitsuteki.etke.host\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"matrix.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=spf1 ip4:$HETZNER_SERVER_IP -all\",\"zone_id\":\"$HETZNER_ZONE_ID\"},{\"name\":\"_dmarc.matrix.higenjitsuteki\",\"type\":\"TXT\",\"value\":\"v=DMARC1; p=quarantine;\",\"zone_id\":\"$HETZNER_ZONE_ID\"}]}"
```

