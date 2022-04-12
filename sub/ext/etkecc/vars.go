package etkecc

import (
	"strings"

	"maunium.net/go/mautrix"
)

func (o *order) generateVars() {
	var txt strings.Builder

	// base
	txt.WriteString(o.generateVarsSystem())
	txt.WriteString(o.generateVarsHomeserver())
	txt.WriteString(o.generateVarsSynapse())
	txt.WriteString(o.generateVarsNginx())

	// additional low-level services
	txt.WriteString(o.generateVarsPostgresBackup())
	txt.WriteString(o.generateVarsCorporal())

	// additional services
	txt.WriteString(o.generateVarsCinny())
	txt.WriteString(o.generateVarsDimension())
	txt.WriteString(o.generateVarsDnsmasq())
	txt.WriteString(o.generateVarsElement())
	txt.WriteString(o.generateVarsEtherpad())
	txt.WriteString(o.generateVarsJitsi())
	txt.WriteString(o.generateVarsKuma())
	txt.WriteString(o.generateVarsLanguagetool())
	txt.WriteString(o.generateVarsMiniflux())
	txt.WriteString(o.generateVarsMiounne())
	txt.WriteString(o.generateVarsRadicale())
	txt.WriteString(o.generateVarsStats())
	txt.WriteString(o.generateVarsSynapseAdmin())
	txt.WriteString(o.generateVarsWireguard())

	// bots
	txt.WriteString(o.generateVarsGoneb())
	txt.WriteString(o.generateVarsHonoroit())
	txt.WriteString(o.generateVarsMjolnir())
	txt.WriteString(o.generateVarsReminder())

	// bridges
	txt.WriteString(o.generateVarsDiscord())
	txt.WriteString(o.generateVarsFacebook())
	txt.WriteString(o.generateVarsGooglechat())
	txt.WriteString(o.generateVarsGroupme())
	txt.WriteString(o.generateVarsHeisenbridge())
	txt.WriteString(o.generateVarsInstagram())
	txt.WriteString(o.generateVarsLinkedin())
	txt.WriteString(o.generateVarsSignal())
	txt.WriteString(o.generateVarsSlack())
	txt.WriteString(o.generateVarsSteam())
	txt.WriteString(o.generateVarsTelegram())
	txt.WriteString(o.generateVarsTwitter())
	txt.WriteString(o.generateVarsWebhooks())
	txt.WriteString(o.generateVarsWhatsapp())

	o.files = append(o.files, &mautrix.ReqUploadMedia{
		Content:       strings.NewReader(txt.String()),
		FileName:      "vars.yml",
		ContentType:   "text/yaml",
		ContentLength: int64(txt.Len()),
	})
}

func (o *order) generateVarsSystem() string {
	var txt strings.Builder

	txt.WriteString("\n# system\n")
	txt.WriteString("system_security_autorizedkeys: []\n")

	return txt.String()
}

func (o *order) generateVarsHomeserver() string {
	var txt strings.Builder

	txt.WriteString("\n# homeserver https://matrix." + o.get("domain") + "\n")
	txt.WriteString("matrix_domain: " + o.get("domain") + "\n")
	txt.WriteString("matrix_admin: \"@" + o.get("username") + ":{{ matrix_domain }}\"\n")
	txt.WriteString("matrix_ssl_lets_encrypt_support_email: " + o.get("email") + "\n")
	if o.has("ma1sd") {
		txt.WriteString("matrix_ma1sd_enabled: no\n")
	}
	txt.WriteString("matrix_mailer_enabled: no\n")
	if !o.has("element-web") {
		txt.WriteString("matrix_client_element_enabled: no\n")
	}

	return txt.String()
}

func (o *order) generateVarsPostgresBackup() string {
	var txt strings.Builder

	txt.WriteString("\n# postgres::backups\n")
	txt.WriteString("matrix_postgres_backup_enabled: yes\n")
	txt.WriteString("matrix_postgres_backup_schedule: '@daily'\n")
	txt.WriteString("matrix_postgres_backup_keep_days: 7\n")
	txt.WriteString("matrix_postgres_backup_keep_weeks: 0\n")
	txt.WriteString("matrix_postgres_backup_keep_months: 0\n")

	return txt.String()
}

func (o *order) generateVarsSynapse() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse\n")
	txt.WriteString("matrix_homeserver_implementation: synapse\n")
	txt.WriteString("matrix_synapse_presence_enabled: yes\n")
	txt.WriteString("matrix_synapse_enable_registration: yes\n")
	txt.WriteString("matrix_synapse_enable_group_creation: yes\n")
	txt.WriteString("matrix_synapse_max_upload_size_mb: 1024\n")
	txt.WriteString("matrix_synapse_tmp_directory_size_mb: \"{{ matrix_synapse_max_upload_size_mb * 2 }}\"\n")

	txt.WriteString("\n# synapse::federation\n")
	txt.WriteString("matrix_synapse_allow_public_rooms_over_federation: yes\n")
	txt.WriteString("matrix_synapse_allow_public_rooms_without_auth: yes\n")

	txt.WriteString("\n# synapse::custom\n")
	txt.WriteString("matrix_synapse_configuration_extension_yaml: |\n")
	txt.WriteString("  registration_requires_token: yes\n")
	txt.WriteString("  disable_msisdn_registration: yes\n")
	txt.WriteString("  allow_device_name_lookup_over_federation: no\n")

	txt.WriteString("\n# synapse::privacy\n")
	txt.WriteString("matrix_synapse_user_ips_max_age: 5m\n")
	txt.WriteString("matrix_synapse_redaction_retention_period: 5m\n")

	txt.WriteString("\n# synapse::extensions::shared_secret_auth\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_enabled: yes\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_shared_secret: " + o.pwgen() + "\n")

	if o.has("synapse-simple-antispam") {
		txt.WriteString("\n# synapse::extensions::simple-antispam\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_synapse_simple_antispam_enabled: yes\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_synapse_simple_antispam_config_blocked_homeservers: []\n")
	}

	if o.has("mjolnir") {
		txt.WriteString("\n# synapse::extensions::spam_checker_mjolnir\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_mjolnir_antispam_enabled: yes\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_mjolnir_antispam_config_block_invites: no\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_mjolnir_antispam_config_block_messages: no\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_mjolnir_antispam_config_block_usernames: no\n")
		txt.WriteString("matrix_synapse_ext_spam_checker_mjolnir_antispam_config_ban_lists: []\n")
	}

	if o.has("matrix-corporal") {
		txt.WriteString("\n# synapse::extensions::rest_auth\n")
		txt.WriteString("matrix_synapse_ext_password_provider_rest_auth_enabled: yes\n")
		txt.WriteString("matrix_synapse_ext_password_provider_rest_auth_endpoint: \"http://matrix-corporal:41080/_matrix/corporal\"\n")
	}

	if o.has("synapse-workers") {
		txt.WriteString("\n# synapse::workers\n")
		txt.WriteString("matrix_synapse_workers_enabled: yes\n")
		txt.WriteString("matrix_synapse_workers_preset: one-of-each\n")
	}

	if o.has("smtp-relay") {
		txt.WriteString(o.generateVarsSynapseMailer())
	}

	txt.WriteString(o.generateVarsSynapseCredentials())

	return txt.String()
}

func (o *order) generateVarsSynapseMailer() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse::mailer\n")
	txt.WriteString("matrix_synapse_email_enabled: yes\n")
	txt.WriteString("matrix_synapse_email_smtp_host: smtp.migadu.com\n")
	txt.WriteString("matrix_synapse_email_smtp_port: 587\n")
	txt.WriteString("matrix_synapse_email_smtp_user: \"matrix@{{ matrix_domain }}\"\n")
	txt.WriteString("matrix_synapse_email_smtp_pass: " + o.pwgen() + "\n")
	txt.WriteString("matrix_synapse_email_notif_from: \"Matrix <matrix@{{ matrix_domain }}>\"\n")
	txt.WriteString("matrix_synapse_email_client_base_url: \"https://{{ matrix_server_fqn_element }}\"\n")
	txt.WriteString("matrix_synapse_email_invite_client_location: \"https://{{ matrix_server_fqn_element }}\"\n")

	return txt.String()
}

func (o *order) generateVarsSynapseCredentials() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse::credentials\n")
	txt.WriteString("matrix_synapse_macaroon_secret_key: " + o.pwgen() + "\n")
	txt.WriteString("matrix_postgres_connection_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_synapse_password_config_pepper: " + o.pwgen() + "\n")
	txt.WriteString("matrix_coturn_turn_static_auth_secret: " + o.pwgen() + "\n")
	txt.WriteString("matrix_homeserver_generic_secret_key: \"{{ matrix_synapse_macaroon_secret_key }}\"\n")

	return txt.String()
}

func (o *order) generateVarsCorporal() string {
	if !o.has("matrix-corporal") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# matrix corporal\n")
	txt.WriteString("matrix_corporal_enabled: yes\n")
	txt.WriteString("matrix_corporal_http_api_enabled: yes\n")
	txt.WriteString("matrix_corporal_http_api_auth_token: " + o.password("matrix-corporal api") + "\n")
	txt.WriteString("matrix_corporal_corporal_user_id_local_part: \"corporal\" # password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_corporal_policy_provider_config: |\n")
	txt.WriteString("  # TODO\n")

	return txt.String()
}

func (o *order) generateVarsSynapseAdmin() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse-admin https://matrix." + o.get("domain") + "/synapse-admin\n")
	txt.WriteString("matrix_synapse_admin_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsNginx() string {
	var txt strings.Builder

	txt.WriteString("\n# nginx proxy\n")
	txt.WriteString("matrix_nginx_proxy_access_log_enabled: no\n")
	txt.WriteString("matrix_nginx_proxy_base_domain_serving_enabled: " + o.get("serve_base_domain") + "\n")
	txt.WriteString("matrix_nginx_proxy_base_domain_homepage_enabled: no\n")
	txt.WriteString(o.generateVarsNginxCustom())

	return txt.String()
}

func (o *order) generateVarsNginxCustom() string {
	var has bool
	var txt strings.Builder

	txt.WriteString("matrix_ssl_additional_domains_to_obtain_certificates_for:\n")
	for _, custom := range customlist {
		if o.has(custom) {
			has = true
			txt.WriteString("- \"{{ matrix_server_fqn_" + custom + " }}\"\n")
		}
	}

	if has {
		return txt.String()
	}

	return ""
}

func (o *order) generateVarsCinny() string {
	if !o.has("cinny") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# cinny https://cinny." + o.get("domain") + "\n")
	txt.WriteString("matrix_client_cinny_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsDimension() string {
	if !o.has("dimension") {
		return ""
	}
	var txt strings.Builder
	password := o.pwgen()

	txt.WriteString("\n# dimension https://dimension." + o.get("domain") + "\n")
	txt.WriteString("matrix_dimension_enabled: yes\n")
	txt.WriteString("matrix_dimension_access_token: TODO # password: " + password + "\n")
	txt.WriteString("matrix_dimension_admins:\n")
	txt.WriteString("  - \"{{ matrix_admin }}\"")
	txt.WriteString("# TODO\n")
	txt.WriteString("# matrix-synapse-register-user dimension " + password + " 0\n")
	txt.WriteString("# curl -X POST -H 'Content-Type: application/json' -d '{\"identifier\": { \"type\": \"m.id.user\", \"user\": \"dimension\" },\"password\": \"" + password + "\", \"type\": \"m.login.password\"}' 'https://matrix." + o.get("domain") + "/_matrix/client/r0/login'\n")

	return txt.String()
}

func (o *order) generateVarsEtherpad() string {
	if !o.has("etherpad") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# etherpad\n")
	txt.WriteString("matrix_etherpad_enabled: yes\n")
	txt.WriteString("matrix_etherpad_admin_username: " + o.get("username") + "\n")
	txt.WriteString("matrix_etherpad_admin_password: " + o.password("etherpad admin") + "\n")

	return txt.String()
}

func (o *order) generateVarsElement() string {
	if !o.has("element-web") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# element https://element." + o.get("domain") + "\n")
	txt.WriteString("matrix_client_element_enabled: yes\n")
	txt.WriteString("matrix_client_element_default_theme: dark\n")

	return txt.String()
}

func (o *order) generateVarsJitsi() string {
	if !o.has("jitsi") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# jitsi https://jitsi." + o.get("domain") + "\n")
	txt.WriteString("matrix_jitsi_enabled: yes\n")
	txt.WriteString("# matrix_jitsi_enable_auth: yes\n")
	txt.WriteString("# matrix_jitsi_enable_guests: yes\n")
	txt.WriteString("matrix_jitsi_jvb_auth_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_jitsi_jibri_xmpp_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_jitsi_jibri_recorder_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_jitsi_jicofo_auth_password: " + o.pwgen() + "\n")
	txt.WriteString("# matrix_jitsi_prosody_auth_internal_accounts:\n")
	txt.WriteString("#  - username: " + o.get("username") + "\n")
	txt.WriteString("#    password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) generateVarsStats() string {
	if !o.has("stats") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# stats https://stats." + o.get("domain") + "\n")
	txt.WriteString("matrix_grafana_enabled: yes\n")
	txt.WriteString("matrix_prometheus_enabled: yes\n")
	txt.WriteString("matrix_grafana_anonymous_access: no\n")
	txt.WriteString("matrix_prometheus_node_exporter_enabled: yes\n")
	txt.WriteString("matrix_grafana_default_admin_user: " + o.get("username") + "\n")
	txt.WriteString("matrix_grafana_default_admin_password: " + o.password("grafana") + "\n")

	return txt.String()
}

func (o *order) generateVarsDnsmasq() string {
	if !o.has("dnsmasq") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# dnsmasq\n")
	txt.WriteString("custom_dnsmasq_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsKuma() string {
	if !o.has("kuma") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# uptime-kuma https://kuma." + o.get("domain") + "\n")
	txt.WriteString("custom_kuma_enabled: yes\n")
	txt.WriteString("matrix_server_fqn_kuma: \"kuma.{{ matrix_domain }}\"\n")

	return txt.String()
}

func (o *order) generateVarsLanguagetool() string {
	if !o.has("languagetool") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# uptime-languagetool https://languagetool." + o.get("domain") + "\n")
	txt.WriteString("custom_languagetool_enabled: yes\n")
	txt.WriteString("custom_languagetool_ngrams_enabled: yes # WARNING: requires a LOT of storage\n")
	txt.WriteString("matrix_server_fqn_languagetool: \"languagetool.{{ matrix_domain }}\"\n")

	return txt.String()
}

func (o *order) generateVarsRadicale() string {
	if !o.has("radicale") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# radicale https://radicale." + o.get("domain") + "\n")
	txt.WriteString("custom_radicale_enabled: yes\n")
	txt.WriteString("matrix_server_fqn_radicale: \"radicale.{{ matrix_domain }}\"\n")
	txt.WriteString("custom_radicale_htpasswd: \"" + o.get("username") + ":TODO\"\n")
	txt.WriteString("# TODO: htpasswd -nb " + o.get("username") + " " + o.password("radicale") + "\n")

	return txt.String()
}

func (o *order) generateVarsMiniflux() string {
	if !o.has("miniflux") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# miniflux https://miniflux." + o.get("domain") + "\n")
	txt.WriteString("custom_miniflux_enabled: yes\n")
	txt.WriteString("matrix_server_fqn_miniflux: \"miniflux.{{ matrix_domain }}\"\n")
	txt.WriteString("# TODO:\n")
	txt.WriteString("# matrix-postgres-cli-non-interactive -c \"CREATE USER miniflux WITH PASSWORD '" + o.pwgen() + "';\"\n")
	txt.WriteString("# matrix-postgres-cli-non-interactive -c \"CREATE DATABASE miniflux; GRANT ALL PRIVILEGES ON DATABASE miniflux to miniflux;\"\n")
	txt.WriteString("# docker exec -it custom-miniflux /usr/bin/miniflux -create-admin\n")

	return txt.String()
}

func (o *order) generateVarsMiounne() string {
	if !o.has("miounne") {
		return ""
	}
	var txt strings.Builder
	password := o.pwgen()

	txt.WriteString("\n# miounne https://miounne." + o.get("domain") + "\n")
	txt.WriteString("custom_miounne_enabled: yes\n")
	txt.WriteString("custom_miounne_matrix_user_login: miounne\n")
	txt.WriteString("custom_miounne_matrix_user_password: " + password + "\n")

	txt.WriteString("# TODO: only for registration\n")
	txt.WriteString("custom_miounne_registration_url: https://matrix." + o.get("domain") + "/matrix-registration\n")
	txt.WriteString("custom_miounne_matrix_registration_room: TODO\n")
	txt.WriteString("custom_miounne_matrix_registration_secret: TODO\n")

	txt.WriteString("# TODO: only for BMC\n")
	txt.WriteString("custom_miounne_bmc_token: TODO\n")
	txt.WriteString("custom_miounne_bmc_room: TODO\n")
	txt.WriteString("custom_miounne_bmc_notify_extras: 1\n")
	txt.WriteString("custom_miounne_bmc_notify_members: 1\n")
	txt.WriteString("custom_miounne_bmc_notify_supporters: 1\n")

	txt.WriteString("# TODO: only for forms\n")
	txt.WriteString("matrix_server_fqn_miounne: \"miounne.{{ matrix_domain }}\"\n")
	txt.WriteString("matrix_nginx_proxy_proxy_miounne_hostname: \"{{ matrix_server_fqn_miounne }}\"\n")
	txt.WriteString("custom_miounne_spam_emails: []\n")
	txt.WriteString("custom_miounne_spam_hosts: []\n")
	txt.WriteString("custom_miounne_forms: []\n")

	txt.WriteString("# TODO: matrix-synapse-register-user miounne " + password + " 0\n")

	return txt.String()
}

func (o *order) generateVarsWireguard() string {
	if !o.has("wireguard") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# wireguard\n")
	txt.WriteString("custom_wireguard_enabled: yes\n")
	txt.WriteString("custom_wireguard_clients: [] # TODO\n")

	return txt.String()
}

func (o *order) generateVarsGoneb() string {
	if !o.has("go-neb") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::goneb\n")
	txt.WriteString("matrix_bot_go_neb_enabled: yes\n")
	txt.WriteString("matrix_bot_go_neb_clients:\n")
	txt.WriteString("- UserID: \"@goneb:{{ matrix_domain }}\"\n")
	txt.WriteString("  AccessToken: \"TODO\" # password: " + o.pwgen() + "\n")
	txt.WriteString("  DeviceID: server\n")
	txt.WriteString("  HomeserverURL: \"{{ matrix_homeserver_container_url }}\"\n")
	txt.WriteString("  Sync: yes\n")
	txt.WriteString("  AutoJoinRooms: yes\n")
	txt.WriteString("  DisplayName: \"GoNEB\"\n")
	txt.WriteString("  AcceptVerificationFromUsers: [\":{{ matrix_domain }}\"]\n")
	txt.WriteString("matrix_bot_go_neb_services: [] # TODO\n")

	return txt.String()
}

func (o *order) generateVarsHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	pass := o.pwgen()
	txt.WriteString("\n# bots::honoroit\n")
	txt.WriteString("matrix_bot_honoroit_enabled: yes\n")
	txt.WriteString("matrix_bot_honoroit_password: " + pass + "\n")
	txt.WriteString("matrix_bot_honoroit_roomid: TODO\n")
	txt.WriteString("# TODO: matrix-synapse-register-user honoroit " + pass + " 0\n")

	return txt.String()
}

func (o *order) generateVarsMjolnir() string {
	if !o.has("mjolnir") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::mjolnir\n")
	txt.WriteString("matrix_bot_mjolnir_enabled: yes\n")
	txt.WriteString("matrix_bot_mjolnir_access_token: \"TODO\" # password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_mjolnir_management_room: TODO\n")
	txt.WriteString("matrix_bot_mjolnir_configuration_extension_yaml: |\n")
	txt.WriteString("  recordIgnoredInvites: true\n")

	return txt.String()
}

func (o *order) generateVarsReminder() string {
	if !o.has("reminder-bot") {
		return ""
	}
	var txt strings.Builder
	password := o.pwgen()

	txt.WriteString("\n# bots::reminder\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_enabled: yes\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_reminders_timezone: TODO\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_matrix_user_id_localpart: reminder\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_matrix_user_password: " + password + "\n")
	txt.WriteString("# TODO: matrix-synapse-register-user reminder " + password + " 0\n")

	return txt.String()
}

func (o *order) generateVarsDiscord() string {
	if !o.has("discord") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::discord\n")
	txt.WriteString("matrix_mx_puppet_discord_enabled: yes\n")
	txt.WriteString("matrix_mx_puppet_discord_configuration_extension_yaml: |\n")
	txt.WriteString("  presence:\n")
	txt.WriteString("    enabled: no\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    console: warn\n")

	return txt.String()
}

func (o *order) generateVarsGroupme() string {
	if !o.has("groupme") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::groupme\n")
	txt.WriteString("matrix_mx_puppet_groupme_enabled: yes\n")
	txt.WriteString("matrix_mx_puppet_groupme_configuration_extension_yaml: |\n")
	txt.WriteString("  presence:\n")
	txt.WriteString("    enabled: no\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    console: warn\n")

	return txt.String()
}

func (o *order) generateVarsSlack() string {
	if !o.has("slack") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::slack\n")
	txt.WriteString("matrix_mx_puppet_slack_enabled: yes\n")
	txt.WriteString("matrix_mx_puppet_slack_configuration_extension_yaml: |\n")
	txt.WriteString("  presence:\n")
	txt.WriteString("    enabled: no\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    console: warn\n")

	return txt.String()
}

func (o *order) generateVarsSteam() string {
	if !o.has("steam") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::steam\n")
	txt.WriteString("matrix_mx_puppet_steam_enabled: yes\n")
	txt.WriteString("matrix_mx_puppet_steam_configuration_extension_yaml: |\n")
	txt.WriteString("  presence:\n")
	txt.WriteString("    enabled: no\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    console: warn\n")

	return txt.String()
}

func (o *order) generateVarsFacebook() string {
	if !o.has("facebook") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::facebook\n")
	txt.WriteString("matrix_mautrix_facebook_enabled: yes\n")
	txt.WriteString("matrix_mautrix_facebook_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_facebook_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      paho:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsGooglechat() string {
	if !o.has("googlechat") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::googlechat\n")
	txt.WriteString("matrix_mautrix_googlechat_enabled: yes\n")
	txt.WriteString("matrix_mautrix_googlechat_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_googlechat_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      hangups:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsHeisenbridge() string {
	if !o.has("heisenbridge") {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("\n# bridges::irc (heisenbridge)\n")
	txt.WriteString("matrix_heisenbridge_enabled: yes\n")
	txt.WriteString("matrix_heisenbridge_identd_enabled: yes\n")
	txt.WriteString("matrix_heisenbridge_owner: \"{{ matrix_admin }}\"\n")

	return txt.String()
}

func (o *order) generateVarsInstagram() string {
	if !o.has("instagram") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::instagram\n")
	txt.WriteString("matrix_mautrix_instagram_enabled: yes\n")
	txt.WriteString("matrix_mautrix_instagram_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_instagram_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      mauigpapi:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      paho:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsLinkedin() string {
	if !o.has("linkedin") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::linkedin\n")
	txt.WriteString("matrix_beeper_linkedin_enabled: yes\n")
	txt.WriteString("matrix_beeper_linkedin_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_beeper_linkedin_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      paho:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsSignal() string {
	if !o.has("signal") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::signal\n")
	txt.WriteString("matrix_mautrix_signal_enabled: yes\n")
	txt.WriteString("matrix_mautrix_signal_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_signal_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsTelegram() string {
	if !o.has("telegram") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::telegram\n")
	txt.WriteString("matrix_mautrix_telegram_enabled: yes\n")
	txt.WriteString("matrix_mautrix_telegram_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    delivery_error_reports: yes\n")
	txt.WriteString("    delivery_receipts: no\n")
	txt.WriteString("    max_initial_member_sync: 10\n")
	txt.WriteString("    sync_channel_members: no\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_telegram_homeserver_domain }}\": full\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      telethon:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsTwitter() string {
	if !o.has("twitter") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::twitter\n")
	txt.WriteString("matrix_mautrix_twitter_enabled: yes\n")
	txt.WriteString("matrix_mautrix_twitter_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    delivery_error_reports: yes\n")
	txt.WriteString("    delivery_receipts: no\n")
	txt.WriteString("    presence: no\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_twitter_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    loggers:\n")
	txt.WriteString("      mau:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      aiohttp:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("    root:\n")
	txt.WriteString("      level: WARNING\n")
	txt.WriteString("      handlers: [console]\n")

	return txt.String()
}

func (o *order) generateVarsWebhooks() string {
	if !o.has("webhooks") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::webhooks\n")
	txt.WriteString("matrix_appservice_webhooks_enabled: yes\n")
	txt.WriteString("matrix_appservice_webhooks_api_secret: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) generateVarsWhatsapp() string {
	if !o.has("whatsapp") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::whatsapp\n")
	txt.WriteString("matrix_mautrix_whatsapp_enabled: yes\n")
	txt.WriteString("matrix_mautrix_whatsapp_configuration_extension_yaml: |\n")
	txt.WriteString("  bridge:\n")
	txt.WriteString("    permissions:\n")
	txt.WriteString("      \"{{ matrix_mautrix_whatsapp_homeserver_domain }}\": user\n")
	txt.WriteString("      \"{{ matrix_admin }}\": admin\n")
	txt.WriteString("  logging:\n")
	txt.WriteString("    print_level: warn\n")

	return txt.String()
}
