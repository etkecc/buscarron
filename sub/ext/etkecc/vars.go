package etkecc

import (
	"strings"

	"maunium.net/go/mautrix"
)

func (o *order) generateVars() {
	var txt strings.Builder

	// base
	txt.WriteString(o.generateVarsAll())
	txt.WriteString(o.generateVarsPostgres())
	txt.WriteString(o.generateVarsHomeserver())
	txt.WriteString(o.generateVarsUsers())
	txt.WriteString(o.generateVarsSynapse())
	txt.WriteString(o.generateVarsNginx())

	// additional low-level services
	txt.WriteString(o.generateVarsBorgBackup())
	txt.WriteString(o.generateVarsPostgresBackup())
	txt.WriteString(o.generateVarsSygnal())
	txt.WriteString(o.generateVarsNtfy())

	// additional services
	txt.WriteString(o.generateVarsCinny())
	txt.WriteString(o.generateVarsDimension())
	txt.WriteString(o.generateVarsElement())
	txt.WriteString(o.generateVarsEtherpad())
	txt.WriteString(o.generateVarsHydrogen())
	txt.WriteString(o.generateVarsJitsi())
	txt.WriteString(o.generateVarsStats())
	txt.WriteString(o.generateVarsSynapseAdmin())

	// bots
	txt.WriteString(o.generateVarsBuscarron())
	txt.WriteString(o.generateVarsHonoroit())
	txt.WriteString(o.generateVarsPostmoogle())
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
	txt.WriteString(o.generateVarsSkype())
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

func (o *order) generateVarsAll() string {
	var txt strings.Builder

	txt.WriteString("### all:start\n")
	txt.WriteString("### all:end\n")

	return txt.String()
}

func (o *order) generateVarsPostgres() string {
	var txt strings.Builder

	txt.WriteString("\n# postgres\n")
	txt.WriteString("devture_postgres_connection_password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) generateVarsHomeserver() string {
	var txt strings.Builder

	txt.WriteString("\n# homeserver https://matrix." + o.get("domain") + "\n")
	txt.WriteString("matrix_domain: " + o.get("domain") + "\n")
	txt.WriteString("matrix_admin: \"@" + o.get("username") + ":{{ matrix_domain }}\"\n")
	txt.WriteString("matrix_ssl_lets_encrypt_support_email: " + o.get("email") + "\n")
	txt.WriteString("matrix_mailer_enabled: no\n")
	if !o.has("element-web") {
		txt.WriteString("matrix_client_element_enabled: no\n")
	}

	return txt.String()
}

func (o *order) generateVarsUsers() string {
	var txt strings.Builder

	txt.WriteString("\n# initial users\n")
	txt.WriteString("matrix_user_creator_users_additional:\n")
	txt.WriteString(" - username: " + o.get("username") + "\n")
	txt.WriteString("   initial_password: " + o.password("matrix") + "\n")
	txt.WriteString("   initial_type: admin\n")

	return txt.String()
}

func (o *order) generateVarsSygnal() string {
	if !o.has("sygnal") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# sygnal https://sygnal." + o.get("domain") + "\n")
	txt.WriteString("matrix_sygnal_enabled: yes\n")
	txt.WriteString("matrix_sygnal_apps:\n")
	txt.WriteString("  TODO:\n")
	txt.WriteString("    type: gcm\n")
	txt.WriteString("    api_key: TODO\n")
	txt.WriteString("matrix_sygnal_configuration_extension_yaml:\n")
	txt.WriteString("  log:\n")
	txt.WriteString("    setup:\n")
	txt.WriteString("      root:\n")
	txt.WriteString("        level: WARNING\n")
	txt.WriteString("      loggers:\n")
	txt.WriteString("        sygnal:\n")
	txt.WriteString("          level: WARNING\n")
	txt.WriteString("        sygnal.access:\n")
	txt.WriteString("          level: WARNING\n")

	return txt.String()
}

func (o *order) generateVarsNtfy() string {
	if !o.has("ntfy") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# ntfy https://ntfy." + o.get("domain") + "\n")
	txt.WriteString("matrix_ntfy_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsPostgresBackup() string {
	if o.has("borg") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# postgres::backups\n")
	txt.WriteString("devture_postgres_backup_enabled: yes\n")
	txt.WriteString("devture_postgres_backup_schedule: '@daily'\n")
	txt.WriteString("devture_postgres_backup_keep_days: 7\n")
	txt.WriteString("devture_postgres_backup_keep_weeks: 0\n")
	txt.WriteString("devture_postgres_backup_keep_months: 0\n")

	return txt.String()
}

func (o *order) generateVarsBorgBackup() string {
	if !o.has("borg") {
		return ""
	}
	var txt strings.Builder
	pub, priv := o.keygen()

	txt.WriteString("\n# borg\n")
	txt.WriteString("matrix_backup_borg_enabled: yes\n")
	txt.WriteString("matrix_backup_borg_location_repositories: [] # TODO\n")
	txt.WriteString("matrix_backup_borg_storage_encryption_passphrase: " + o.pwgen() + "\n")
	txt.WriteString("matrix_backup_borg_ssh_key_private: |\n")
	for _, line := range strings.Split(priv, "\n") {
		txt.WriteString("  " + line + "\n")
	}
	txt.WriteString("# TODO: " + pub + "\n")

	return txt.String()
}

func (o *order) generateVarsSynapse() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse\n")
	txt.WriteString("matrix_homeserver_implementation: synapse\n")
	txt.WriteString("matrix_synapse_presence_enabled: yes\n")
	txt.WriteString("matrix_synapse_enable_registration: yes\n")
	txt.WriteString("matrix_synapse_registration_requires_token: yes\n")
	txt.WriteString("matrix_synapse_max_upload_size_mb: 1024\n")
	txt.WriteString("matrix_synapse_media_retention_remote_media_lifetime: 30d\n")
	txt.WriteString("matrix_synapse_tmp_directory_size_mb: \"{{ matrix_synapse_max_upload_size_mb * 2 }}\"\n")

	txt.WriteString("\n# synapse::federation\n")
	txt.WriteString("matrix_synapse_allow_public_rooms_over_federation: yes\n")
	txt.WriteString("matrix_synapse_allow_public_rooms_without_auth: yes\n")

	txt.WriteString("\n# synapse::custom\n")
	txt.WriteString("matrix_synapse_configuration_extension_yaml: |\n")
	txt.WriteString("  disable_msisdn_registration: yes\n")
	if o.has("sso") {
		txt.WriteString("  oidc_providers:\n")
		txt.WriteString("  - idp_id: todo\n")
		txt.WriteString("    idp_name: Todo\n")
		txt.WriteString("    idp_brand: \"todo\"\n")
		txt.WriteString("    issuer: \"https://TODO/\"\n")
		txt.WriteString("    client_id: \"TODO\"\n")
		txt.WriteString("    client_secret: \"TODO\"\n")
		txt.WriteString("    scopes: [\"openid\", \"profile\"]\n")
		txt.WriteString("    user_mapping_provider:\n")
		txt.WriteString("      config:\n")
		txt.WriteString("        localpart_template: \"{% raw %}{{ user.given_name|lower }}{% endraw %}\"\n")
		txt.WriteString("        display_name_template: \"{% raw %}{{ user.name }}{% endraw %}\"\n")
	}

	txt.WriteString("\n# synapse::privacy\n")
	txt.WriteString("matrix_synapse_user_ips_max_age: 5m\n")
	txt.WriteString("matrix_synapse_redaction_retention_period: 5m\n")

	txt.WriteString("\n# synapse::extensions::shared_secret_auth\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_enabled: yes\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_shared_secret: " + o.pwgen() + "\n")

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

	return txt.String()
}

func (o *order) generateVarsSynapseCredentials() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse::credentials\n")
	txt.WriteString("matrix_synapse_macaroon_secret_key: " + o.pwgen() + "\n")
	txt.WriteString("matrix_synapse_password_config_pepper: " + o.pwgen() + "\n")
	txt.WriteString("matrix_coturn_turn_static_auth_secret: " + o.pwgen() + "\n")
	txt.WriteString("matrix_homeserver_generic_secret_key: \"{{ matrix_synapse_macaroon_secret_key }}\"\n")

	return txt.String()
}

func (o *order) generateVarsSynapseAdmin() string {
	if !o.has("synapse-admin") {
		return ""
	}

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
	txt.WriteString(o.generateVarsNginxWebsite())

	return txt.String()
}

func (o *order) generateVarsNginxWebsite() string {
	if !o.has("nginx-proxy-website") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# nginx proxy website\n")
	txt.WriteString("matrix_nginx_proxy_website_enabled: yes\n")
	txt.WriteString("matrix_nginx_proxy_website_repo: TODO\n")
	txt.WriteString("matrix_nginx_proxy_website_command: TODO\n")
	txt.WriteString("matrix_nginx_proxy_website_dist: \"public\"\n")

	return txt.String()
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
	if o.has("dimension") {
		txt.WriteString("matrix_etherpad_mode: dimension\n")
	}
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

func (o *order) generateVarsHydrogen() string {
	if !o.has("hydrogen") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# hydrogen https://hydrogen." + o.get("domain") + "\n")
	txt.WriteString("matrix_client_hydrogen_enabled: yes\n")
	txt.WriteString("matrix_client_hydrogen_default_theme: dark\n")

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
	txt.WriteString("matrix_prometheus_node_exporter_process_extra_arguments:\n")
	txt.WriteString("  - \"--collector.systemd\"\n")
	txt.WriteString("matrix_grafana_default_admin_user: " + o.get("username") + "\n")
	txt.WriteString("matrix_grafana_default_admin_password: " + o.password("grafana") + "\n")

	return txt.String()
}

func (o *order) generateVarsBuscarron() string {
	if !o.has("buscarron") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::buscarron\n")
	txt.WriteString("matrix_bot_buscarron_enabled: yes\n")
	txt.WriteString("matrix_bot_buscarron_login: buscarron\n")
	txt.WriteString("matrix_bot_buscarron_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_buscarron_forms: [] # TODO\n")

	return txt.String()
}

func (o *order) generateVarsHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::honoroit\n")
	txt.WriteString("matrix_bot_honoroit_enabled: yes\n")
	txt.WriteString("matrix_bot_honoroit_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_honoroit_roomid: TODO\n")

	return txt.String()
}

func (o *order) generateVarsPostmoogle() string {
	if !o.has("postmoogle") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::postmoogle\n")
	txt.WriteString("matrix_bot_postmoogle_enabled: yes\n")
	txt.WriteString("matrix_bot_postmoogle_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_postmoogle_data_secret: " + o.pwgen(32) + "\n")

	return txt.String()
}

func (o *order) generateVarsReminder() string {
	if !o.has("reminder-bot") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::reminder\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_enabled: yes\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_reminders_timezone: TODO\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_matrix_user_id_localpart: reminder\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_matrix_user_password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) generateVarsDiscord() string {
	if !o.has("discord") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::discord\n")
	txt.WriteString("matrix_mautrix_discord_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsGroupme() string {
	if !o.has("groupme") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::groupme\n")
	txt.WriteString("matrix_mx_puppet_groupme_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsSlack() string {
	if !o.has("slack") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::slack\n")
	txt.WriteString("matrix_mx_puppet_slack_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsSteam() string {
	if !o.has("steam") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::steam\n")
	txt.WriteString("matrix_mx_puppet_steam_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsFacebook() string {
	if !o.has("facebook") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::facebook\n")
	txt.WriteString("matrix_mautrix_facebook_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsGooglechat() string {
	if !o.has("googlechat") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::googlechat\n")
	txt.WriteString("matrix_mautrix_googlechat_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsHeisenbridge() string {
	if !o.has("irc") {
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

	return txt.String()
}

func (o *order) generateVarsLinkedin() string {
	if !o.has("linkedin") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::linkedin\n")
	txt.WriteString("matrix_beeper_linkedin_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsSignal() string {
	if !o.has("signal") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::signal\n")
	txt.WriteString("matrix_mautrix_signal_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsSkype() string {
	if !o.has("skype") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::skype\n")
	txt.WriteString("matrix_go_skype_bridge_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsTelegram() string {
	if !o.has("telegram") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::telegram\n")
	txt.WriteString("matrix_mautrix_telegram_enabled: yes\n")
	txt.WriteString("matrix_mautrix_telegram_api_id: TODO\n")
	txt.WriteString("matrix_mautrix_telegram_api_hash: TODO\n")

	return txt.String()
}

func (o *order) generateVarsTwitter() string {
	if !o.has("twitter") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::twitter\n")
	txt.WriteString("matrix_mautrix_twitter_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsWebhooks() string {
	if !o.has("webhooks") && !o.has("hookshot") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::hookshot\n")
	txt.WriteString("matrix_hookshot_enabled: yes\n")

	return txt.String()
}

func (o *order) generateVarsWhatsapp() string {
	if !o.has("whatsapp") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::whatsapp\n")
	txt.WriteString("matrix_mautrix_whatsapp_enabled: yes\n")

	return txt.String()
}
