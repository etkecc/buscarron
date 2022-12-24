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
		"q_dimension":           "are you sure you want it? It's a self-hosted integration manager. You will get an integration manager by default with any Element client app. Please check https://github.com/turt2live/matrix-dimension before deciding",
		"q_nginx-proxy-website": "to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).",
		"q_sso":                 "You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)",
		"q_sygnal":              "are you sure you want it? It's a push gateway, usable only for matrix client app developers, so you can't use it if you don't develop your mobile matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)",
		"q_borg":                "please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.",
		"q_jitsi":               "are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.",
		"q_turnkey_ssh":         "SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.",
		"q_byos_ssh":            "Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.\nAdd our SSH keys (https://etke.cc/ssh.key) to your server, send us your server's IP address, the username (with permissions to call sudo), and password (if set).",

		"dns_add_entries":    "DNS - please, add the following entries",
		"ps_automatic_email": "PS: this is an automated email. Please, reply to it with answers to the questions above (if any). An operator (human) will proceed with your answers",
		"matrix_server_on":   "Matrix server on",
		"auth_instructions":  "auth instructions",
		"in_case_of_issues":  "in case of any issues",
		"credentials":        "credentials",
		"links":              "links",
		"bridges":            "bridges",
		"bots":               "bots",
		"payment":            "payment",
		"buy_setup":          "Please, [buy the Setup item]",
		"join_the":           "join the",
		"membership_on":      "membership on",
		"and":                "and",
		"note_pwyw":          "> **NOTE**: all prices are based on the [Pay What You Want](https://en.wikipedia.org/wiki/Pay_what_you_want) model.",

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
	"de": {
		"intro": "Hallo zusammen,\nWir haben Ihre Bestellung erhalten und haben einige Fragen vor der Einrichtung.",

		"q_delegation":          "Wir sehen, dass Sie etwas auf Ihrer Basisdomain haben. In diesem Fall sollten Sie die folgenden HTTPS-Weiterleitungen (HTTP 301) hinzufügen",
		"q_reminder-bot":        "Was ist Ihre Zeitzone (IANA)? Wie America/Chicago, Asia/Seoul, oder Europe/Berlin",
		"q_buscarron":           "Sind Sie sicher, dass Sie es wollen? Es ist ein Bot, der Webformulare (HTML/HTTP POST) empfängt und sie an (verschlüsselte) Matrixräume sendet. Bitte prüfen Sie https://gitlab.com/etke.cc/buscarron, bevor Sie sich entscheiden.",
		"q_honoroit":            "Sind Sie sicher, dass Sie es wollen? Es ist ein Helpdesk-Bot mit Ende-Zu-Ende-Verschlüsselungsunterstützung. Bitte prüfen Sie https://gitlab.com/etke.cc/honoroit, bevor Sie sich entscheiden.",
		"q_telegram":            "gehen Sie bitte auf https://my.telegram.org/apps und erstellen Sie eine neue App. Teilen Sie mir die API-ID und den Hash mit",
		"q_smtp-relay":          "Bitte wählen Sie einen geeigneten E-Mail-Anbieter (große Anbieter wie Gmail oder Outlook sperren Sie für automatisierte E-Mails, daher müssen Sie einen Dienst finden, der den Versand von Verifizierungs-E-Mails erlaubt. Optional bieten wir einen solchen Dienst an). Bitte senden Sie uns einen SMTP Host, einen SMTP STARTTLS Port, ein SMTP Login, ein SMTP Passwort und eine SMTP E-Mail (normalerweise sind Login und E-Mail dasselbe, aber das hängt vom Anbieter ab).",
		"q_stats":               "Sind Sie sicher, dass Sie das wollen? Cloud-Anbieter bieten in der Regel ein Dashboard mit Server-Statistiken, warum also nicht stattdessen dieses Dashboard verwenden? Ein Prometheus+Grafana-Stack bietet einige interne Matrix-Statistiken (z. B. die Anzahl der Ereignisse), ist aber zu viel des Guten, wenn Sie nur die Serverauslastung sehen wollen.",
		"q_dimension":           "Sind Sie sicher, dass Sie das wollen? Es ist ein selbst gehosteter Integrationsmanager. Sie erhalten standardmäßig einen Integrationsmanager mit jeder Element-Client-Anwendung. Bitte prüfen Sie https://github.com/turt2live/matrix-dimension bevor Sie sich entscheiden.",
		"q_nginx-proxy-website": "Um eine statische Website einzurichten, müssen Sie Ihre Basisdomäne (den @ DNS eintrag) auf die IP-Adresse des Matrix-Servers verweisen, und der Quellcode der Website muss in einem öffentlichen git repo verfügbar sein. Sind Sie sicher, dass Sie das wollen? Wenn ja, geben Sie bitte die URL des Website-Repositorys, den Befehl (-s) zum Erstellen der Website und den Ordner an, in dem die Build-Distribution gespeichert ist (normalerweise public oder dist).",
		"q_sso":                 "Sie haben nicht erwähnt, welchen OIDC/OAuth2-Anbieter Sie integrieren möchten. Hier ist eine Liste der gängigen Anbieter - https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs. Bitte senden Sie uns die Informationen, die für die Konfiguration erforderlich sind (in der Regel sind es Anbietername, issuer, client_id, client_secret, aber das hängt vom Anbieter ab)",
		"q_sygnal":              "Sind Sie sicher, dass Sie es wollen? Es ist ein Push Gateway, das nur für Matrix-Client-App-Entwickler nutzbar ist. Sie können es also nicht verwenden, wenn Sie Ihre mobile Matrix-App nicht entwickeln. Wenn Sie es hinzufügen möchten, geben Sie bitte die folgenden Informationen an: App ID(-s) (z.B. org.matrix.app), FCM api key, und/oder APNS certificate (falls verwendet)",
		"q_borg":                "geben Sie bitte die gewünschte Repository-URL an (user@host:repo). Wir werden einen SSH-Schlüssel und eine Verschlüsselungspassphrase auf Ihrem Server generieren. Wir senden Ihnen den öffentlichen Teil des generierten SSH-Schlüssels zu. Sie müssen diesen SSH-Schlüssel zu Ihrem Provider hinzufügen.",
		"q_jitsi":               "Sind Sie sicher, dass Sie es wollen? Sie erhalten die jitsi-Integration standardmäßig mit einer öffentlichen Instanz. Das von uns angebotene jitsi ist eine selbst gehostete Version. Denken Sie daran, dass Jitsi die Anforderungen an die Rechenleistung deutlich erhöht.",
		"q_turnkey_ssh":         "SSH: Sie bestellen einen gehosteten/verwalteten Server. Wir werden den Server in Ihrem Namen einrichten und verwalten. Dennoch können Sie vollen SSH-Zugang zu diesem Server erhalten. **Wenn** Sie SSH-Zugang zu diesem Server wünschen, senden Sie uns Ihren öffentlichen SSH-Schlüssel und eine Liste der IP-Adressen, von denen aus Sie auf den Server zugreifen möchten.",
		"q_byos_ssh":            "Server: Bitte erstellen Sie einen x86/amd64 VPS mit einer beliebigen Debian-basierten Distribution. Minimale komfortable Konfiguration für einen einfachen Matrix-Server: 1vCPU, 2GB RAM.\n Fügen Sie unsere SSH-Schlüssel (https://etke.cc/ssh.key) zu Ihrem Server hinzu, senden Sie uns die IP-Adresse Ihres Servers, den Benutzernamen (mit der Berechtigung, sudo aufzurufen) und das Passwort (falls festgelegt).",

		"dns_add_entries":    "DNS - bitte fügen Sie die folgenden Einträge hinzu",
		"ps_automatic_email": "PS: Dies ist eine automatisierte E-Mail. Bitte antworten Sie auf diese E-Mail und beantworten Sie die oben gestellten Fragen (falls vorhanden). Ein Operator (Mensch) wird Ihre Antworten bearbeiten.",
		"matrix_server_on":   "Matrix-Server an",
		"auth_instructions":  "Authentifizierungsanweisungen",
		"in_case_of_issues":  "im Falle irgendwelcher Probleme",
		"credentials":        "Anmeldeinformationen",
		"links":              "links",
		"bridges":            "bridges",
		"bots":               "bots",
		"payment":            "Bezahlung",
		"buy_setup":          "Bitte, [kaufen Sie das Einrichtungselement]",
		"join_the":           "Beitritt zur",
		"membership_on":      "Mitgliedschaft auf",
		"and":                "und",
		"note_pwyw":          "> **HINWEIS**: Alle Preise basieren auf dem [Pay What You Want]-Modell (https://de.wikipedia.org/wiki/Pay_what_you_want).",

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
