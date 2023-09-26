package etkecc

const i18nDefault = "en"

var i18n map[string]map[string]string = map[string]map[string]string{
	"en": {
		"intro":                 "Hi there,\nWe got your order and have some questions before the setup.",
		"q_delegation":          "We see that you have something on your base domain. In that case, you should add the following HTTPS redirects (HTTP 301)",
		"q_delegation_details":  "To learn more about why these redirects are necessary and what the connection between the base domain (DOMAIN) and the Matrix domain (matrix.DOMAIN) is, read the following guide: " + link("etke.cc/help/faq#why-do-i-need-well-known-redirects-on-the-base-domain"),
		"q_reminder-bot":        "What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin",
		"q_buscarron":           "are you sure you want it? It's a bot that receives web forms (HTML/HTTP POST) and send them to (encrypted) Matrix rooms. Please, check https://gitlab.com/etke.cc/buscarron before deciding.",
		"q_honoroit":            "are you sure you want it? It's a helpdesk bot with e2e encryption support. Please, check " + link("gitlab.com/etke.cc/honoroit") + " before deciding.",
		"q_telegram":            "please, go to " + link("https://my.telegram.org/apps") + " and create a new app. Share the API ID and Hash with us",
		"q_smtp-relay":          "please, select a suitable email provider (big providers like Gmail or Outlook will ban you for automated emails, so you need to find a service that allows sending of verification emails. Optionally, we provide such service). Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email (usually login and email are the same thing, but that depends on the provider).",
		"q_stats":               "are you sure you want it? Cloud providers usually provide a dashboard with server stats, so why not use that dashboard instead? A Prometheus+Grafana stack provides some internal Matrix stats (like count of events), but it's overkill if you just want to see server utilization.",
		"q_nginx-proxy-website": "to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).",
		"q_sso":                 "You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - " + link("github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs") + ". Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)",
		"q_sygnal":              "are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)",
		"q_borg":                "please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.",
		"q_jitsi":               "are you sure you want it? You will get jitsi integration by default with a public instance. The jitsi item we offer is a self-hosted version. Keep in mind that jitsi significantly increases compute power requirements.",
		"q_turnkey_ssh":         "SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.",
		"q_byos_ssh":            "Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a basic Matrix server: 1vCPU, 2GB RAM.\nAdd our SSH keys (" + link("etke.cc/ssh.key") + ") to your server, open the required ports (" + link("etke.cc/help/faq#what-ports-should-be-open") + ") send us your server's IP address, the username (with permissions to call sudo), and password (if set).",
		"dns_add_entries":       "Please, add the following DNS entries",
		"ps_automatic_email":    "PS: this is an automated email. Please, reply to it with answers to the questions above (if any). An operator (human) will proceed with your answers",
		"hello":                 "Hello!",
		"server_is_ready":       "Your Matrix server is ready, here is your onboarding list:",
		"happy_matrixing":       "Happy Matrixing!",
		"ps_matrixrooms":        "PS: To help with your first steps in the world of Matrix, we've built the " + link("MatrixRooms.info") + " search engine. You can use it to discover rooms over the Matrix Federation and find amazing communities you wish to join!",
		"matrix_server_on":      "Matrix server on",
		"auth_instructions":     "auth instructions",
		"in_case_of_issues":     "In case of any issues",
		"credentials":           "Credentials",
		"links":                 "Links",
		"bridges":               "Bridges",
		"bots":                  "Bots",
		"payment":               "Payment",
		"membership_on":         "membership on",
		"and":                   "and",
		"note_pwyw":             "> **NOTE**: all prices are based on the [Pay What You Want](https://en.wikipedia.org/wiki/Pay_what_you_want) model.",
		"steps_after_setup":     "Steps after the setup",
		"as_honoroit_1":         "Create a matrix room (encryption supported) and invite the honoroit user into it",
		"as_honoroit_2":         "Send the room id to @support:etke.cc",
		"as_buscarron_1":        "Create (encrypted) Matrix room(-s) for the forms you want to have and invite the buscarron user to all of them",
		"as_buscarron_2":        "Send the list of (form name, Matrix room id, redirect URL after submission) to @support:etke.cc",
	},
}

// t is translation func
func t(lang, key string) string {
	if _, ok := i18n[lang]; !ok {
		lang = i18nDefault //nolint:ineffassign
	}

	lang = i18nDefault //nolint:ineffassign // first phase of disabling transations
	v, ok := i18n[lang][key]
	if !ok {
		return key
	}

	return v
}
