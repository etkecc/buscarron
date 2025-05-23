package etkecc

import (
	"context"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
)

func (o *order) generateDelegationInstructions(ctx context.Context) string {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateDelegationInstructions")
	defer span.Finish()
	log := o.logger(ctx)
	log.Info().Msg("generating delegation instructions")

	if o.get("serve_base_domain") == "yes" {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("\nWe see that you have something on your base domain.\n")
	txt.WriteString("**If** that's a domain registrar's (parking) page and/or you intend to serve base domain (" + o.domain + ") from the matrix server, just add the `@` DNS record pointing to the server IP and tell us about that.\n")
	txt.WriteString("**If** that's your website and/or you intend to serve base domain from some other server, you should add the following HTTPS redirects (HTTP 301):\n")
	txt.WriteString("* " + link(o.domain+"/.well-known/matrix/server") + " -> " + link("matrix."+o.domain+"/.well-known/matrix/server") + "\n")
	txt.WriteString("* " + link(o.domain+"/.well-known/matrix/client") + " -> " + link("matrix."+o.domain+"/.well-known/matrix/client") + "\n")
	txt.WriteString("* " + link(o.domain+"/.well-known/matrix/support") + " -> " + link("matrix."+o.domain+"/.well-known/matrix/support") + "\n")
	txt.WriteString("\nTo learn more about why these redirects are necessary and what the connection between the base domain (" + o.domain + ") and the Matrix domain (matrix." + o.domain + ") is, read the following guide: " + link("etke.cc/order/status#delegation-redirects") + "\n\n")

	log.Info().Msg("delegation instructions have been generated")
	return txt.String()
}

func (o *order) generateQuestions(ctx context.Context) (text string, count int) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateQuestions")
	defer span.Finish()
	log := o.logger(ctx)
	log.Info().Msg("generating questions")

	count = 0
	var txt strings.Builder
	if q := o.generateQuestionsReminderBot(); q != "" {
		count++
		txt.WriteString(q)
	}
	if q := o.generateQuestionsTelegramBridge(); q != "" {
		count++
		txt.WriteString(q)
	}
	if q := o.generateQuestionsServices(); q != "" {
		count++
		txt.WriteString(q)
	}
	if q := o.generateQuestionsType(); q != "" {
		count++
		txt.WriteString(q)
	}

	log.Info().Msg("questions have been generated")
	return txt.String(), count
}

func (o *order) generateQuestionsReminderBot() string {
	var txt strings.Builder

	if o.has("reminder-bot") && !o.has("reminder-bot-tz") {
		txt.WriteString("Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin. ")
		txt.WriteString("[Full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsTelegramBridge() string {
	var txt strings.Builder

	if o.has("telegram") && !(o.has("telegram-api-id") && o.has("telegram-api-hash")) {
		txt.WriteString("Telegram: please, go to " + link("https://my.telegram.org/apps") + " and create a new app. ")
		txt.WriteString("Share the API ID and Hash with us\n\n")
	}

	return txt.String()
}

//nolint:gocognit // TODO
func (o *order) generateQuestionsServices() string {
	var txt strings.Builder

	if o.has("smtp-relay") && !o.has("service-email") && !(o.has("smtp-relay-host") && o.has("smtp-relay-port") && o.has("smtp-relay-login") && o.has("smtp-relay-password") && o.has("smtp-relay-email")) {
		txt.WriteString("SMTP relay: please, select a suitable email provider ")
		txt.WriteString("(big providers like Gmail or Outlook will ban you for automated emails, ")
		txt.WriteString("so you need to find a service that allows sending of verification emails. Optionally, we provide such service). ")
		txt.WriteString("Please, send us an SMTP host, SMTP STARTTLS port, SMTP login, SMTP password, and SMTP email ")
		txt.WriteString("(usually login and email are the same thing, but that depends on the provider)\n\n")
	}
	if (o.has("nginx-proxy-website") && o.get("serve_base_domain") == "yes") && !(o.has("nginx-proxy-website-command") && o.has("nginx-proxy-website-repo") && o.has("nginx-proxy-website-dist")) {
		txt.WriteString("Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Supported generators: hugo, jekyll, plain html (no generator). ")
		txt.WriteString("Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).\n\n")
	}
	if (o.has("sso") || o.has("synapse-sso")) && !(o.has("sso-client-id") && o.has("sso-client-secret") && o.has("sso-issuer") && o.has("sso-idp-brand") && o.has("sso-idp-id") && o.has("sso-idp-name")) {
		txt.WriteString("SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - ")
		txt.WriteString(link("github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs") + ". ")
		txt.WriteString("Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)\n\n")
	}
	if o.has("sygnal") && !(o.has("sygnal-app-id") && o.has("sygnal-gcm-apikey")) {
		txt.WriteString("Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, ")
		txt.WriteString("so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, ")
		txt.WriteString("provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)\n\n")
	}
	if o.has("borg") && !o.has("borg-repository") {
		txt.WriteString("BorgBackup: please, provide the desired repository url (user@host:repo). ")
		txt.WriteString("We will generate an SSH key and encryption passphrase on your server. ")
		txt.WriteString("We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsType() string {
	if o.hosting != "" && !o.has("ssh-client-key-disabled") && !o.has("ssh-client-key") {
		return "SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.\n\n"
	}

	if o.hosting == "" && !(o.has("ssh-host") && o.has("ssh-port") && o.has("ssh-user")) {
		return "Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 2vCPU, 2GB RAM.\nOpen the required ports (" + link("etke.cc/order/status/#ports-and-firewalls") + ") send us your server's IP address, the username (with permissions to call sudo), and password (if set).\n\n"
	}
	return ""
}
