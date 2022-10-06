package etkecc

const i18nDefault = "en"

var i18n map[string]map[string]string = map[string]map[string]string{
	"en": {
		"intro": "Hi there,\nWe got your order and have some questions before the setup.",

		"q_delegation":          "We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301)",
		"q_reminder-bot":        "What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin",
		"q_buscarron":           "are you sure you want it? It's a bot that receives web forms (HTML/HTTP POST) and send them to (encrypted) matrix rooms. Please, check https://gitlab.com/etke.cc/buscarron before deciding.",
		"q_honoroit":            "are you sure you want it? It's a helpdesk bot with e2e encryption support. Please, check https://gitlab.com/etke.cc/honoroit before deciding.",
		"q_telegram":            "please, go to https://my.telegram.org/apps and create a new app. Share the API ID and Hash with me",
		"q_smtp-relay":          "please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).",
		"q_stats":               "are you sure you want it? Cloud providers usually provide a dashboard with server stats, so why not use that dashboard instead? A Prometheus+Grafana stack provides some internal matrix stats (like count of events), but it's overkill if you just want to see server utilization.",
		"q_etherpad":            "are you sure you want it? It's a self-hosted collaborative editor, available for usage only with dimension (added a question about it, too). Keep in mind that you will get an etherpad anyway with the default integration manager, but it will be hosted by Element Inc. (developers of the Element client apps).",
		"q_dimension":           "are you sure you want it? It's a self-hosted integration manager. You will get an integration manager by default with any Element client app. Please check https://github.com/turt2live/matrix-dimension before deciding",
		"q_nginx-proxy-website": "to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).",
		"q_sso":                 "You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)",
		"q_sygnal":              "are you sure you want it? It's a push gateway, usable only for matrix client app developers, so you can't use it if you don't develop your mobile matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)",
		"q_borg":                "please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.",
		"q_jitsi":               "are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.",
		"q_turnkey_ssh":         "SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.",
		"q_byos_ssh":            "Server: please, create a VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.\nAdd our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).",

		"only_with_subscription": "only with subscription or turnkey",
		"dns_add_entries":        "DNS - please, add the following entries",
		"ps_automatic_email":     "PS: this is an automated email. Please, reply to it with answers to the questions above (if any). An operator (human) will proceed with your answers",
		"matrix_server_on":       "Matrix server on",
		"auth_instructions":      "auth instructions",
		"in_case_of_issues":      "in case of any issues",
		"credentials":            "credentials",
		"links":                  "links",
		"bridges":                "bridges",
		"bots":                   "bots",
		"payment":                "payment",
		"buy_setup":              "Please, [buy the Setup item]",
		"join_the":               "join the",
		"membership_on":          "membership on",
		"and":                    "and",
		"note_pwyw":              "> **NOTE**: all prices are based on the [Pay What You Want](https://en.wikipedia.org/wiki/Pay_what_you_want) model.",

		"steps_after_setup": "steps after the setup",
		"as_etherpad_1":     "Open integration manager in any element client app",
		"as_etherpad_2":     "Click on the sprocket icon and go to the Widgets tab",
		"as_etherpad_3":     "Open the Etherpad widget configuration and replace `scalar.vector.im` with",
		"as_honoroit_1":     "Create a matrix room (encryption supported) and invite the honoroit user into it",
		"as_honoroit_2":     "Send the room id to the @support:etke.cc",
		"as_honoroit_3":     "if you want to change honoroit's messages, send the texts you want to use to @support:etke.cc (you can use https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go as a reference)",
		"as_buscarron_1":    "Create (encrypted) matrix room(-s) for the forms you want to have and invite the buscarron user to all of them",
		"as_buscarron_2":    "Send the list of (form name, matrix room id, redirect URL after submission) to the @support:etke.cc",
	},
	"ru": {
		"intro": "Привет,\nМы получили Ваш заказ и хотели задать несколько вопросов перед установкой.",

		"q_delegation":          "Мы видим, что у Вас что-то настроено на основном домене. В этом случае Вам нужно будет добавить HTTPS редиректы (HTTP 301)",
		"q_reminder-bot":        "Какой у Вас часовой пояс (в формате IANA)? Например, America/Chicago, Asia/Seoul или Europe/Berlin",
		"q_buscarron":           "Вы точно хотите это? Это бот, который принимает отправку веб форм (HTML/HTTP POST) и отправляет их в (зашифрованные) Matrix комнаты. Пожалуйста, посмотрите на https://gitlab.com/etke.cc/buscarron и решите.",
		"q_honoroit":            "Вы точно хотите это? Это хелпдеск бот с поддержкой e2e шифрования. Пожалуйста, посмотрите на https://gitlab.com/etke.cc/honoroit и решите.",
		"q_telegram":            "пожалуйста, сходите на https://my.telegram.org/apps и создайте новое приложение. Отправьте нам API ID и Hash",
		"q_smtp-relay":          "пожалуйста, выберите подходящего почтового провайдера (большие провайдеры вроде Gmail или Outlook забанят Вас за автоматическую отправку писем, так что Вам нужно найти сервис, который разрешает отправку верификационных имейлов. К примеру, мы предоставляем такой сервис). Пожалуйста, отправьте нам SMTP хост, SMTP STARTTLS порт, SMTP логин, SMTP пароль и SMTP email (обычно логин и email это то же самое, но все зависит от провайдера).",
		"q_stats":               "Вы точно хотите их? Любой хостинг провайдер предоставляет статистику сервера в панели управления, почему бы не использовать ее? Стек Prometheus+Grafana позволяет увидеть внутреннюю статистику matrix (например, количество ивентов), но это как стрелять из пушки по воробьям, если Вам нужно просто посмотреть использование ресурсов на сервере.",
		"q_etherpad":            "Вы точно хотите это? Это редактор текста для совместной работы на Вашем сервере, работает только в связке с dimension (добавили вопрос и о нем). Etherpad и так доступен Вам по умолчанию в менеджере интеграций в любом приложении Element (хостится на серверах Element Inc., разработчики приложений Element).",
		"q_dimension":           "Вы точно хотите это? Это менеджер интеграций на Вашем сервере. По умолчанию у Вас и так будет доступен стандартный менеджер интеграций. Пожалуйста, посмотрите на https://github.com/turt2live/matrix-dimension и решите",
		"q_nginx-proxy-website": "чтобы задеплоить статичный вебсайт, Вам придется настроить свой базовый домен на сервер matrix (@ DNS запись), а исходники самого вебсайта должны быть доступны в публичном git репозитории. Вы точно хотите этого? Если да, пожалуйста, отправьте нам адрес git репозитория Вашего статичного сайта, список команд для сборки и в какой директории будет собран артефакт (обычно это public или dist).",
		"q_sso":                 "Мы не получили информацию о том, какого OIDC/OAuth2 провайдера Вы хотите интегрировать, так что вот список популярных провайдеров - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Пожалуйста, отправьте нам требуемую информацию для настройки (обычно это название провайдера, issuer, client_id и client_secret, но все зависит от выбранного провайдера)",
		"q_sygnal":              "Вы точно хотите это? Это push gateway, нужен только разработчикам мобильных приложений matrix, так что он будет бесполезен, если Вы не разрабатываете собственное мобильное приложение для matrix. Если Вы все же хотите его добавить, пожалуйста, отправьте следующую информацию: ID приложений (например org.matrix.app), API ключ FCM и/или сертификат APNS (если используется)",
		"q_borg":                "пожалуйста, предоставьте желаемый адрес репозитория (user@host:repo). Мы сгенерируем ssh ключ и пароль шифрования на стороне Вашего сервера и отправим Вам публичную часть ключа. Вам нужно будет добавить этот ключ на стороне Вашего borg провайдера.",
		"q_jitsi":               "Вы точно хотите это? По умолчанию у Вас будет интегрирован публичный сервер jitsi, а мы предлагаем то же самое установить на Ваш сервер. Имейте в виду, что jitsi на Вашем сервере сильно увеличивает требования к конфигурации сервера.",
		"q_turnkey_ssh":         "SSH: Вы заказываете хостинг, мы подготовим Ваш сервер и будем его обслуживать. **Если** Вы хотите иметь полный SSH доступ к Вашему серверу, пришлите нам свой публичный SSH ключ и список IP адресов, с которых Вы хотите подключаться к серверу по SSH.",
		"q_byos_ssh":            "Сервер: пожалуйста, создайте VPS с любым дистрибутивом на базе Debian. Минимальная комфортная конфигурация для базового сервера matrix: 1vCPU, 2ГБ RAM.\nДобавьте наши SSH ключи (https://etke.cc/ssh.key) на свой сервер, отправьте нам IP адрес сервера, имя пользователя (с доступом к вызову sudo) и пароль (если установлен).",

		"only_with_subscription": "только с подпиской или хостингом",
		"dns_add_entries":        "DNS - пожалуйста, добавьте следующие записи",
		"ps_automatic_email":     "PS: это автоматическое письмо. Пожалуйста, отправьте свои ответы на вопросы выше прямо ответом на него. Оператор (человек) подключится после получения ответов",
		"matrix_server_on":       "сервер matrix на",
		"auth_instructions":      "инструкции по авторизации",
		"in_case_of_issues":      "в случае проблем",
		"credentials":            "доступы",
		"links":                  "ссылки",
		"bridges":                "интеграции",
		"bots":                   "боты",
		"payment":                "оплата",
		"buy_setup":              "Пожалуйста, [оплатите установку]",
		"join_the":               "оформите",
		"membership_on":          "подписку на",
		"and":                    "и",
		"note_pwyw":              "> **ВАЖНО**: все цены основаны на модели  [Плати Сколько Хочешь](https://ru.wikipedia.org/wiki/%D0%9F%D0%BB%D0%B0%D1%82%D0%B8_%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%BE_%D1%85%D0%BE%D1%87%D0%B5%D1%88%D1%8C).",

		"steps_after_setup": "после установки",
		"as_etherpad_1":     "Откройте менеджер интеграций в любом приложении Element",
		"as_etherpad_2":     "Нажмите на иконку шестеренки и перейдите в таб Widgets",
		"as_etherpad_3":     "Откройте конфигурацию виджета Etherpad и замените `scalar.vector.im` на",
		"as_honoroit_1":     "Создайте matrix комнату (шифрование поддерживается) и пригласите пользователя honoroit в нее",
		"as_honoroit_2":     "Отправьте ID этой комнаты на @support:etke.cc",
		"as_honoroit_3":     "Если Вы хотите изменить текст сообщений, отправьте желаемый вариант на @support:etke.cc (используйте https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go в качесте справки)",
		"as_buscarron_1":    "Создайте (зашифрованные) комнаты matrix для всех форм, которые Вы собираетесь сделать и пригласите пользователя buscarron в каждую из них",
		"as_buscarron_2":    "Отправьте список (название формы, id комнаты matrix, URL для перенаправления пользователя после отправки) на @support:etke.cc",
	},
	"de": {
		"intro": "Hallo zusammen,\nWir haben Ihre Bestellung erhalten und haben einige Fragen vor der Einrichtung.",

		"q_delegation":          "Wir sehen, dass Sie etwas auf Ihrer Basisdomain haben. In diesem Fall sollten Sie die folgenden HTTPS-Weiterleitungen (HTTP 301) hinzufügen",
		"q_reminder-bot":        "Was ist Ihre Zeitzone (IANA)? Wie America/Chicago, Asia/Seoul, oder Europe/Berlin",
		"q_buscarron":           "Sind Sie sicher, dass Sie es wollen? Es ist ein Bot, der Webformulare (HTML/HTTP POST) empfängt und sie an (verschlüsselte) Matrixräume sendet. Bitte prüfen Sie https://gitlab.com/etke.cc/buscarron, bevor Sie sich entscheiden.",
		"q_honoroit":            "Sind Sie sicher, dass Sie es wollen? Es ist ein Helpdesk-Bot mit Ende-Zu-Ende-Verschlüsselungsunterstützung. Bitte prüfen Sie https://gitlab.com/etke.cc/honoroit, bevor Sie sich entscheiden.",
		"q_telegram":            "gehen Sie bitte auf https://my.telegram.org/apps und erstellen Sie eine neue App. Teilen Sie mir die API-ID und den Hash mit",
		"q_smtp-relay":          "Bitte wählen Sie einen geeigneten E-Mail-Anbieter (große Anbieter wie Gmail oder Outlook sperren Sie für automatisierte E-Mails, daher müssen Sie einen Dienst finden, der den Versand von Verifizierungs-E-Mails erlaubt. Optional bieten wir einen solchen Dienst an). Bitte senden Sie uns einen SMTP Host, einen SMTP STARTTLS Port, ein SMTP Login, ein SMTP Passwort und eine SMTP E-Mail (normalerweise sind Login und E-Mail dasselbe, aber das hängt vom Anbieter ab).",
		"q_stats":               "Sind Sie sicher, dass Sie das wollen? Cloud-Anbieter bieten in der Regel ein Dashboard mit Server-Statistiken, warum also nicht stattdessen dieses Dashboard verwenden? Ein Prometheus+Grafana-Stack bietet einige interne Matrix-Statistiken (z. B. die Anzahl der Ereignisse), ist aber zu viel des Guten, wenn Sie nur die Serverauslastung sehen wollen.",
		"q_etherpad":            "Sind Sie sicher, dass Sie das wollen? Es handelt sich um einen selbst gehosteten kollaborativen Editor, der nur in Verbindung mit Dimension verwendet werden kann (auch dazu wurde eine Frage gestellt). Denken Sie daran, dass Sie mit dem Standard-Integrationsmanager ohnehin ein Etherpad erhalten, das jedoch von Element Inc. gehostet wird (Entwickler der Element-Clientanwendungen).",
		"q_dimension":           "Sind Sie sicher, dass Sie das wollen? Es ist ein selbst gehosteter Integrationsmanager. Sie erhalten standardmäßig einen Integrationsmanager mit jeder Element-Client-Anwendung. Bitte prüfen Sie https://github.com/turt2live/matrix-dimension bevor Sie sich entscheiden.",
		"q_nginx-proxy-website": "Um eine statische Website einzurichten, müssen Sie Ihre Basisdomäne (den @ DNS eintrag) auf die IP-Adresse des Matrix-Servers verweisen, und der Quellcode der Website muss in einem öffentlichen git repo verfügbar sein. Sind Sie sicher, dass Sie das wollen? Wenn ja, geben Sie bitte die URL des Website-Repositorys, den Befehl (-s) zum Erstellen der Website und den Ordner an, in dem die Build-Distribution gespeichert ist (normalerweise public oder dist).",
		"q_sso":                 "Sie haben nicht erwähnt, welchen OIDC/OAuth2-Anbieter Sie integrieren möchten. Hier ist eine Liste der gängigen Anbieter - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Bitte senden Sie uns die Informationen, die für die Konfiguration erforderlich sind (in der Regel sind es Anbietername, issuer, client_id, client_secret, aber das hängt vom Anbieter ab)",
		"q_sygnal":              "Sind Sie sicher, dass Sie es wollen? Es ist ein Push Gateway, das nur für Matrix-Client-App-Entwickler nutzbar ist. Sie können es also nicht verwenden, wenn Sie Ihre mobile Matrix-App nicht entwickeln. Wenn Sie es hinzufügen möchten, geben Sie bitte die folgenden Informationen an: App ID(-s) (z.B. org.matrix.app), FCM api key, und/oder APNS certificate (falls verwendet)",
		"q_borg":                "geben Sie bitte die gewünschte Repository-URL an (user@host:repo). Wir werden einen SSH-Schlüssel und eine Verschlüsselungspassphrase auf Ihrem Server generieren. Wir senden Ihnen den öffentlichen Teil des generierten SSH-Schlüssels zu. Sie müssen diesen SSH-Schlüssel zu Ihrem Provider hinzufügen.",
		"q_jitsi":               "Sind Sie sicher, dass Sie es wollen? Sie erhalten die jitsi-Integration standardmäßig mit einer öffentlichen Instanz. Das von uns angebotene jitsi ist eine selbst gehostete Version. Denken Sie daran, dass Jitsi die Anforderungen an die Rechenleistung deutlich erhöht.",
		"q_turnkey_ssh":         "SSH: Sie bestellen einen gehosteten/verwalteten Server. Wir werden den Server in Ihrem Namen einrichten und verwalten. Dennoch können Sie vollen SSH-Zugang zu diesem Server erhalten. **Wenn** Sie SSH-Zugang zu diesem Server wünschen, senden Sie uns Ihren öffentlichen SSH-Schlüssel und eine Liste der IP-Adressen, von denen aus Sie auf den Server zugreifen möchten.",
		"q_byos_ssh":            "Server: Bitte erstellen Sie einen VPS mit einer beliebigen Debian-basierten Distribution. Minimale komfortable Konfiguration für einen einfachen Matrix-Server: 1vCPU, 2GB RAM.\n Fügen Sie unsere SSH-Schlüssel (https://etke.cc/ssh.key) zu Ihrem Server hinzu, senden Sie uns die IP-Adresse Ihres Servers, den Benutzernamen (mit der Berechtigung, sudo aufzurufen) und das Passwort (falls festgelegt).",

		"only_with_subscription": "nur mit Abonnement oder Komplettpaket",
		"dns_add_entries":        "DNS - bitte fügen Sie die folgenden Einträge hinzu",
		"ps_automatic_email":     "PS: Dies ist eine automatisierte E-Mail. Bitte antworten Sie auf diese E-Mail und beantworten Sie die oben gestellten Fragen (falls vorhanden). Ein Operator (Mensch) wird Ihre Antworten bearbeiten.",
		"matrix_server_on":       "Matrix-Server an",
		"auth_instructions":      "Authentifizierungsanweisungen",
		"in_case_of_issues":      "im Falle irgendwelcher Probleme",
		"credentials":            "Anmeldeinformationen",
		"links":                  "links",
		"bridges":                "bridges",
		"bots":                   "bots",
		"payment":                "Bezahlung",
		"buy_setup":              "Bitte, [kaufen Sie das Einrichtungselement]",
		"join_the":               "Beitritt zur",
		"membership_on":          "Mitgliedschaft auf",
		"and":                    "und",
		"note_pwyw":              "> **HINWEIS**: Alle Preise basieren auf dem [Pay What You Want]-Modell (https://de.wikipedia.org/wiki/Pay_what_you_want).",

		"steps_after_setup": "Schritte nach dem Einrichten",
		"as_etherpad_1":     "Öffnen Sie den Integrationsmanager in einer beliebigen Element-Client-Anwendung",
		"as_etherpad_2":     "Klicken Sie auf das Zahnradsymbol und gehen Sie auf die Registerkarte Widgets",
		"as_etherpad_3":     "Öffnen Sie die Etherpad-Widget-Konfiguration und ersetzen Sie `scalar.vector.im` durch",
		"as_honoroit_1":     "Erstellen Sie einen Matrixraum (Verschlüsselung unterstützt) und laden Sie den honoroit-Benutzer dazu ein",
		"as_honoroit_2":     "Senden Sie die Raum-ID an @support:etke.cc",
		"as_honoroit_3":     "Wenn Sie die Nachrichten von honoroit ändern möchten, senden Sie die gewünschten Texte an @support:etke.cc (Sie können https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go als Referenz verwenden)",
		"as_buscarron_1":    "Erstellen Sie (verschlüsselte) Matrixräume für die gewünschten Formulare und laden Sie den buscarron-Benutzer zu allen Räumen ein",
		"as_buscarron_2":    "Senden Sie die Liste (Formularname, Matrix-Raum-ID, Umleitungs-URL nach Übermittlung) an @support:etke.cc",
	},
}

// t is translation func
func t(lang, key string) string {
	if _, ok := i18n[lang]; !ok {
		lang = i18nDefault
	}

	v, ok := i18n[lang][key]
	if !ok {
		return key
	}

	return v
}
