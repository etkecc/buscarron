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
	txt.WriteString("* https://" + domain + "/.well-known/matrix/client -> https://matrix." + domain + "/.well-known/matrix/client\n\n")

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

func (o *order) generateQuestionsServices() string {
	var txt strings.Builder

	txt.WriteString(o.generateQuestionsServicesSystem())
	txt.WriteString(o.generateQuestionsServicesNonMatrix())
	txt.WriteString(o.generateQuestionsServicesSubscribers())

	return txt.String()
}

func (o *order) generateQuestionsServicesSystem() string {
	var txt strings.Builder

	if o.has("smtp-relay") && o.get("type") != "turnkey" {
		txt.WriteString("SMTP relay: " + o.t("q_smtp-relay") + "\n\n")
	}
	if o.has("stats") && o.get("type") != "turnkey" {
		txt.WriteString("Prometheus+Grafana: " + o.t("q_stats") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsServicesSubscribers() string {
	var txt strings.Builder
	if o.has("etherpad") && o.get("dimension") == "auto" {
		txt.WriteString("Etherpad (" + o.t("only_with_subscription") + "): " + o.t("q_etherpad") + "\n\n")
	}
	if o.has("dimension") {
		txt.WriteString("Dimension (" + o.t("only_with_subscription") + "): " + o.t("q_dimension") + "\n\n")
	}
	if o.has("nginx-proxy-website") {
		txt.WriteString("Website (" + o.t("only_with_subscription") + "): " + o.t("q_nginx-proxy-website") + "\n\n")
	}
	if o.has("buscarron") {
		txt.WriteString("buscarron (" + o.t("only_with_subscription") + "): " + o.t("q_buscarron") + "\n\n")
	}
	if o.has("sso") {
		txt.WriteString("SSO (" + o.t("only_with_subscription") + "): " + o.t("q_sso") + "\n\n")
	}
	if o.has("sygnal") {
		txt.WriteString("Sygnal (" + o.t("only_with_subscription") + "): " + o.t("q_sygnal") + "\n\n")
	}
	if o.has("borg") {
		txt.WriteString("BorgBackup (" + o.t("only_with_subscription") + "): " + o.t("q_borg") + "\n\n")
	}
	if o.has("email2matrix") {
		txt.WriteString("email2matrix (" + o.t("only_with_subscription") + "): " + o.t("q_email2matrix") + "\n\n")
	}
	if o.has("jitsi") {
		txt.WriteString("Jitsi (" + o.t("only_with_subscription") + "): " + o.t("q_jitsi") + "\n\n")
	}
	if o.has("ma1sd") {
		txt.WriteString("ma1sd (" + o.t("only_with_subscription") + "): " + o.t("q_ma1sd") + "\n\n")
	}
	if o.has("matrix-registration") {
		txt.WriteString("matrix-registration (" + o.t("only_with_subscription") + "): " + o.t("q_matrix-registration") + "\n\n")
	}
	if o.has("miounne") {
		txt.WriteString("Miounne (" + o.t("only_with_subscription") + "): " + o.t("q_miounne") + "\n\n")
	}

	return txt.String()
}

func (o *order) generateQuestionsServicesNonMatrix() string {
	var txt strings.Builder
	if o.has("kuma") {
		txt.WriteString("Uptime Kuma: " + o.t("q_kuma") + "\n\n")
	}
	if o.has("radicale") {
		txt.WriteString("Radicale: " + o.t("q_radicale") + "\n\n")
	}
	if o.has("miniflux") {
		txt.WriteString("Miniflux: " + o.t("q_miniflux") + "\n\n")
	}
	if o.has("languagetool") {
		txt.WriteString("Languagetool: " + o.t("q_languagetool") + "\n\n")
	}
	if o.has("softserve") {
		txt.WriteString("Soft-Serve: " + o.t("q_softserve") + "\n\n")
	}
	txt.WriteString(o.generateQuestionsServicesWireguard())

	return txt.String()
}

func (o *order) generateQuestionsServicesWireguard() string {
	switch {
	case o.has("wireguard") && o.has("dnsmasq"):
		return o.t("q_wireguard_dnsmasq") + "\n\n"
	case o.has("wireguard") && !o.has("dnsmasq"):
		return o.t("q_wireguard") + "\n\n"
	case !o.has("wireguard") && o.has("dnsmasq"):
		return o.t("q_dnsmasq") + "\n\n"
	default:
		return ""
	}
}

func (o *order) generateQuestionsType() string {
	if o.get("type") == "turnkey" {
		return o.t("q_turnkey_ssh") + "\n\n"
	}

	return o.t("q_byos_ssh") + "\n\n"
}
