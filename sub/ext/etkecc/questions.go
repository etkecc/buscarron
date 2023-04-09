package etkecc

import "strings"

func (o *order) generateQuestions() string {
	var txt strings.Builder
	txt.WriteString(o.t("intro") + "\n\n")
	txt.WriteString(o.generateQuestionsDelegation())
	txt.WriteString(o.generateQuestionsBots())
	txt.WriteString(o.generateQuestionsBridges())
	txt.WriteString(o.generateQuestionsServices())

	txt.WriteString(o.generateQuestionsType())
	return txt.String()
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

func (o *order) generateQuestionsBots() string {
	var txt strings.Builder

	if o.has("reminder-bot") {
		txt.WriteString("Reminder bot: " + o.t("q_reminder-bot") + "\n\n")
	}
	if o.has("honoroit") {
		txt.WriteString("Honoroit: " + o.t("q_honoroit") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsBridges() string {
	var txt strings.Builder

	if o.has("telegram") {
		txt.WriteString("Telegram: " + o.t("q_telegram") + "\n\n")
	}

	return txt.String()
}

// nolint:gocognit
func (o *order) generateQuestionsServices() string {
	var txt strings.Builder

	if o.has("smtp-relay") {
		txt.WriteString("SMTP relay: " + o.t("q_smtp-relay") + "\n\n")
	}
	if o.has("stats") && o.get("type") != "turnkey" {
		txt.WriteString("Prometheus+Grafana: " + o.t("q_stats") + "\n\n")
	}
	if o.has("dimension") {
		txt.WriteString("Dimension: " + o.t("q_dimension") + "\n\n")
	}
	if o.has("nginx-proxy-website") {
		txt.WriteString("Website: " + o.t("q_nginx-proxy-website") + "\n\n")
	}
	if o.has("buscarron") {
		txt.WriteString("buscarron: " + o.t("q_buscarron") + "\n\n")
	}
	if o.has("sso") {
		txt.WriteString("SSO: " + o.t("q_sso") + "\n\n")
	}
	if o.has("sygnal") {
		txt.WriteString("Sygnal: " + o.t("q_sygnal") + "\n\n")
	}
	if o.has("borg") {
		txt.WriteString("BorgBackup: " + o.t("q_borg") + "\n\n")
	}
	if o.has("jitsi") {
		txt.WriteString("Jitsi: " + o.t("q_jitsi") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsType() string {
	if o.get("type") == "turnkey" {
		return o.t("q_turnkey_ssh") + "\n\n"
	}

	return o.t("q_byos_ssh") + "\n\n"
}
