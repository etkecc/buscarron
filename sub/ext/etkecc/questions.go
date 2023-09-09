package etkecc

import "strings"

func (o *order) generateQuestions() (string, int) {
	var count int
	var txt strings.Builder
	txt.WriteString(o.t("intro") + "\n\n")
	if q := o.generateQuestionsDelegation(); q != "" {
		count++
		txt.WriteString(q)
	}
	if q := o.generateQuestionsReminderBot(); q != "" {
		count++
		txt.WriteString(q)
	}
	if q := o.generateQuestionsHonoroit(); q != "" {
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
	return txt.String(), count
}

func (o *order) generateQuestionsDelegation() string {
	if o.get("serve_base_domain") == "yes" {
		return ""
	}

	var txt strings.Builder
	domain := o.get("domain")
	txt.WriteString(o.t("q_delegation") + ":\n")
	txt.WriteString("* https://" + domain + "/.well-known/matrix/server -> https://matrix." + domain + "/.well-known/matrix/server\n")
	txt.WriteString("* https://" + domain + "/.well-known/matrix/client -> https://matrix." + domain + "/.well-known/matrix/client\n")
	txt.WriteString(strings.ReplaceAll(o.t("q_delegation_details"), "DOMAIN", domain) + "\n\n")

	return txt.String()
}

func (o *order) generateQuestionsReminderBot() string {
	var txt strings.Builder

	if o.has("reminder-bot") && !o.has("reminder-bot-tz") {
		txt.WriteString("Reminder bot: " + o.t("q_reminder-bot") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsHonoroit() string {
	var txt strings.Builder

	if o.has("honoroit") {
		txt.WriteString("Honoroit: " + o.t("q_honoroit") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsTelegramBridge() string {
	var txt strings.Builder

	if o.has("telegram") && !(o.has("telegram-api-id") && o.has("telegram-api-hash")) {
		txt.WriteString("Telegram: " + o.t("q_telegram") + "\n\n")
	}

	return txt.String()
}

// nolint:gocognit
func (o *order) generateQuestionsServices() string {
	var txt strings.Builder

	if o.has("smtp-relay") && !(o.has("smtp-relay-host") && o.has("smtp-relay-port") && o.has("smtp-relay-login") && o.has("smtp-relay-password") && o.has("smtp-relay-email")) {
		txt.WriteString("SMTP relay: " + o.t("q_smtp-relay") + "\n\n")
	}
	if o.has("stats") && o.get("type") != "turnkey" {
		txt.WriteString("Prometheus+Grafana: " + o.t("q_stats") + "\n\n")
	}
	if o.has("nginx-proxy-website") && !(o.has("nginx-proxy-website-command") && o.has("nginx-proxy-website-dist") && o.has("nginx-proxy-website-dist")) {
		txt.WriteString("Website: " + o.t("q_nginx-proxy-website") + "\n\n")
	}
	if o.has("buscarron") {
		txt.WriteString("buscarron: " + o.t("q_buscarron") + "\n\n")
	}
	if o.has("sso") && !(o.has("sso-client-id") && o.has("sso-client-secret") && o.has("sso-issuer") && o.has("sso-idp-brand") && o.has("sso-idp-id") && o.has("sso-idp-name")) {
		txt.WriteString("SSO: " + o.t("q_sso") + "\n\n")
	}
	if o.has("sygnal") && !(o.has("sygnal-app-id") && o.has("sygnal-gcm-apikey")) {
		txt.WriteString("Sygnal: " + o.t("q_sygnal") + "\n\n")
	}
	if o.has("borg") && !o.has("borg-repository") {
		txt.WriteString("BorgBackup: " + o.t("q_borg") + "\n\n")
	}
	if o.has("jitsi") {
		txt.WriteString("Jitsi: " + o.t("q_jitsi") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsType() string {
	if o.getHostingSize() != "" {
		return o.t("q_turnkey_ssh") + "\n\n"
	}

	if !(o.has("ssh-host") && o.has("ssh-port") && o.has("ssh-user")) {
		return o.t("q_byos_ssh") + "\n\n"
	}
	return ""
}
