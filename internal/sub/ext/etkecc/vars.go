package etkecc

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
	"maunium.net/go/mautrix"
)

func (o *order) vars(ctx context.Context) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.vars")
	defer span.Finish()

	log := o.logger(span.Context())
	log.Info().Msg("generating vars")

	var txt strings.Builder

	// base
	txt.WriteString(o.varsEtke())
	txt.WriteString(o.varsSSH())
	txt.WriteString(o.varsPostgres())
	txt.WriteString(o.varsHomeserver())
	txt.WriteString(o.varsMSC1929())
	txt.WriteString(o.varsUsers())
	txt.WriteString(o.varsSynapse())
	txt.WriteString(o.varsStaticFiles())

	// additional low-level services
	txt.WriteString(o.varsBorgBackup())
	txt.WriteString(o.varsEximRelay())
	txt.WriteString(o.varsMatrixRTC())
	txt.WriteString(o.varsNtfy())
	txt.WriteString(o.varsPostgresBackup())
	txt.WriteString(o.varsSygnal())

	// additional services
	txt.WriteString(o.varsCinny())
	txt.WriteString(o.varsElement())
	txt.WriteString(o.varsElementCall())
	txt.WriteString(o.varsEtherpad())
	txt.WriteString(o.varsFirezone())
	txt.WriteString(o.varsFluffyChat())
	txt.WriteString(o.varsFunkwhale())
	txt.WriteString(o.varsGoToSocial())
	txt.WriteString(o.varsHydrogen())
	txt.WriteString(o.varsJitsi())
	txt.WriteString(o.varsLanguagetool())
	txt.WriteString(o.varsLinkding())
	txt.WriteString(o.varsMiniflux())
	txt.WriteString(o.varsPeertube())
	txt.WriteString(o.varsRadicale())
	txt.WriteString(o.varsSchildiChat())
	txt.WriteString(o.varsStats())
	txt.WriteString(o.varsSynapseAdmin())
	txt.WriteString(o.varsUptimeKuma())
	txt.WriteString(o.varsVaultwarden())

	// bots
	txt.WriteString(o.varsBaibot())
	txt.WriteString(o.varsBuscarron())
	txt.WriteString(o.varsDraupnir())
	txt.WriteString(o.varsHonoroit())
	txt.WriteString(o.varsMaubot())
	txt.WriteString(o.varsReminder())

	// bridges
	txt.WriteString(o.varsBluesky())
	txt.WriteString(o.varsDiscord())
	txt.WriteString(o.varsEmail())
	txt.WriteString(o.varsFacebook())
	txt.WriteString(o.varsGmessages())
	txt.WriteString(o.varsGooglechat())
	txt.WriteString(o.varsHeisenbridge())
	txt.WriteString(o.varsInstagram())
	txt.WriteString(o.varsLinkedin())
	txt.WriteString(o.varsSignal())
	txt.WriteString(o.varsSlack())
	txt.WriteString(o.varsSteam())
	txt.WriteString(o.varsTelegram())
	txt.WriteString(o.varsTwitter())
	txt.WriteString(o.varsWebhooks())
	txt.WriteString(o.varsWechat())
	txt.WriteString(o.varsWhatsapp())

	text := txt.String()
	o.files = append(o.files, &mautrix.ReqUploadMedia{
		Content:       strings.NewReader(text),
		ContentBytes:  []byte(text),
		FileName:      "vars.yml",
		ContentType:   "text/yaml",
		ContentLength: int64(len(text)),
	})
	log.Info().Msg("vars have been generated")
}

func (o *order) varsEtke() string {
	enabledServices := map[string]any{}
	enabledServices["etke_order_email"] = o.get("email")
	if o.has("issue_id") && o.get("issue_id") != "0" {
		enabledServices["etke_order_issue_id"] = o.get("issue_id")
	}

	if o.has("matrix") {
		enabledServices["etke_base_matrix"] = "yes"
	}

	if o.has("service-email") {
		enabledServices["etke_service_email"] = "yes"
	}

	if o.has("service-support") {
		enabledServices["etke_service_support"] = o.get("service-support")
	} else {
		enabledServices["etke_service_support"] = "basic"
	}

	o.varsEtkeDNS(enabledServices)
	o.varsEtkeExternalDNS(enabledServices)
	o.varsEtkeHosting(enabledServices)
	o.varsEtkeServices(enabledServices)

	enabledServices["etke_subscription_email"] = o.get("email")
	enabledServices["etke_subscription_provider"] = "Ko-Fi"
	enabledServices["etke_subscription_confirmed"] = "no"

	keys := make([]string, 0, len(enabledServices))
	for k := range enabledServices {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return o.varsEtkeBuilder(keys, enabledServices)
}

func (o *order) varsEtkeDNS(enabledServices map[string]any) {
	if !o.subdomain {
		return
	}
	var serverIPv4, serverIPv6 string
	if o.hosting != "" {
		serverIPv4 = "SERVER_IP4"
		serverIPv6 = "SERVER_IP6"
	} else {
		serverIPv4 = o.get("ssh-host")
	}

	domain := o.domain
	subdomain := strings.Split(domain, ".")[0]
	suffix := "." + subdomain
	var zoneID string
	for sufix, zone := range domains {
		if strings.HasSuffix(domain, sufix) {
			zoneID = zone
			break
		}
	}
	enabledServices["etke_service_dns_zone"] = zoneID
	enabledServices["etke_service_dns_records"] = o.generateDNSRecords(subdomain, suffix, serverIPv4, serverIPv6)
}

func (o *order) varsEtkeExternalDNS(enabledServices map[string]any) {
	if o.subdomain || o.hosting == "" {
		return
	}
	var serverIPv4, serverIPv6 string
	if o.hosting != "" {
		serverIPv4 = "SERVER_IP4"
		serverIPv6 = "SERVER_IP6"
	} else {
		serverIPv4 = o.get("ssh-host")
	}

	if o.v.NS(o.domain, "cloudflare.com") {
		enabledServices["etke_service_dns_external_proxy"] = "yes"
	}

	if o.get("serve_base_domain") != "yes" {
		enabledServices["etke_service_dns_external_delegation"] = "yes"
	}

	enabledServices["etke_service_dns_external_records"] = o.generateDNSRecords("@", "", serverIPv4, serverIPv6)
}

func (o *order) varsEtkeHosting(enabledServices map[string]any) {
	if o.hosting == "" {
		return
	}

	enabledServices["etke_service_server"] = o.hosting
	enabledServices["etke_service_server_location"] = o.location
	enabledServices["etke_service_server_firewalls"] = defaultFirewall
}

func (o *order) varsEtkeServices(enabledServices map[string]any) {
	for field := range o.data {
		_, ok := enabledServices[field]
		if strings.HasPrefix(field, "etke_service") && !ok {
			enabledServices[field] = "yes"
		}
	}
}

func (o *order) varsEtkeBuilder(keys []string, enabledServices map[string]any) string {
	var txt strings.Builder
	txt.WriteString("# etke services\n")
	for _, service := range keys {
		if valueStr, ok := enabledServices[service].(string); ok {
			txt.WriteString(service)
			txt.WriteString(": " + valueStr + "\n")
		}
		if valueSlice, ok := enabledServices[service].([]string); ok {
			txt.WriteString(service)
			txt.WriteString(":\n")
			for _, value := range valueSlice {
				txt.WriteString("  - " + value + "\n")
			}
		}
	}
	return txt.String()
}

func (o *order) varsSSHFirewall() string {
	var txt strings.Builder
	if o.get("ssh-client-ips") == "N/A" {
		txt.WriteString("system_security_ufw_ssh_allowed: yes\n")
		return txt.String()
	}
	txt.WriteString("system_security_ufw_rules_custom:\n")
	for idx, ip := range strings.Split(o.get("ssh-client-ips"), ",") {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		if !strings.Contains(ip, "/") {
			ip += "/32"
		}

		txt.WriteString(fmt.Sprintf("  - name: customer-ssh-%d\n", idx+1))
		txt.WriteString("    rule: allow\n")
		txt.WriteString("    from: \"" + ip + "\"\n")
		txt.WriteString("    port: \"{{ system_security_ssh_port }}\"\n")
		txt.WriteString("    proto: tcp\n")
	}
	return txt.String()
}

func (o *order) varsSSH() string {
	var txt strings.Builder
	txt.WriteString("\n# ssh\n")
	txt.WriteString("ansible_ssh_private_key_file: \"{{ playbook_dir }}/../../inventory/host_vars/{{ inventory_hostname }}/sshkey.priv\"\n")

	if o.has("ssh-port") && o.get("ssh-port") != "22" {
		txt.WriteString("system_security_ssh_port: \"")
		txt.WriteString(o.get("ssh-port"))
		txt.WriteString("\"\n")
	}

	if o.has("ssh-client-ips") {
		txt.WriteString(o.varsSSHFirewall())
	}

	if o.has("ssh-client-key") {
		keys := strings.Split(o.get("ssh-client-key"), "\n")
		txt.WriteString("system_security_ssh_authorizedkeys_host:\n")
		for _, key := range keys {
			txt.WriteString("  - \"")
			txt.WriteString(key)
			txt.WriteString("\"\n")
		}
	}

	if o.get("ssh-user") == "matrix" {
		txt.WriteString("\n# matrix user\n")
		txt.WriteString("matrix_user_name: matrixserver\n")
		txt.WriteString("matrix_group_name: matrixserver\n")
	}

	return txt.String()
}

func (o *order) varsPostgres() string {
	var txt strings.Builder

	txt.WriteString("\n# postgres\n")
	txt.WriteString("postgres_connection_password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) varsHomeserver() string {
	var txt strings.Builder

	txt.WriteString("\n# homeserver\n")
	txt.WriteString("matrix_domain: \"" + o.domain + "\"\n")
	txt.WriteString("matrix_admin: \"@" + o.get("username") + ":" + o.domain + "\"\n")
	if !o.has("element-web") {
		txt.WriteString("matrix_client_element_enabled: no\n")
	}
	if o.has("bridges-encryption") {
		txt.WriteString("matrix_bridges_encryption_enabled: yes\n")
	}

	return txt.String()
}

func (o *order) varsMSC1929() string {
	var txt strings.Builder

	txt.WriteString("\n# MSC1929 admin contacts\n")
	txt.WriteString("matrix_static_files_file_matrix_support_enabled: yes\n")
	txt.WriteString("matrix_static_files_file_matrix_support_property_m_contacts:\n")
	txt.WriteString("  - matrix_id: \"@" + o.get("username") + ":" + o.domain + "\"\n")
	txt.WriteString("    email_address: \"" + o.get("email") + "\"\n")
	txt.WriteString("    role: m.role.admin\n")

	return txt.String()
}

func (o *order) varsUsers() string {
	var txt strings.Builder

	o.login("matrix", o.get("username"))
	txt.WriteString("\n# initial users\n")
	txt.WriteString("matrix_user_creator_users_additional:\n")
	txt.WriteString(" - username: \"" + o.get("username") + "\"\n")
	txt.WriteString("   initial_password: " + o.password("matrix") + "\n")
	txt.WriteString("   initial_type: admin\n")

	if o.has("gotosocial") {
		login := o.login("gotosocial", strings.ReplaceAll(o.get("username"), ".", "_"))
		txt.WriteString("gotosocial_users_additional:\n")
		txt.WriteString(" - username: \"" + login + "\"\n")
		txt.WriteString("   initial_email: \"" + o.get("email") + "\"\n")
		txt.WriteString("   initial_password: " + o.password("gotosocial") + "\n")
		txt.WriteString("   initial_type: admin\n")
	}

	if o.has("funkwhale") {
		login := o.login("funkwhale", strings.ReplaceAll(o.get("username"), ".", "_"))
		txt.WriteString("funkwhale_users_additional:\n")
		txt.WriteString(" - username: \"" + login + "\"\n")
		txt.WriteString("   initial_email: \"" + o.get("email") + "\"\n")
		txt.WriteString("   initial_password: " + o.password("funkwhale") + "\n")
		txt.WriteString("   initial_type: admin\n")
	}

	return txt.String()
}

func (o *order) varsSygnal() string {
	if !o.has("sygnal") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# sygnal\n")
	txt.WriteString("matrix_sygnal_enabled: yes\n")
	txt.WriteString("matrix_sygnal_apps:\n")
	txt.WriteString("  " + o.get("sygnal-app-id") + ":\n")
	txt.WriteString("    type: gcm\n")
	txt.WriteString("    api_key: " + o.get("sygnal-gcm-apikey") + "\n")
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

func (o *order) varsNtfy() string {
	if !o.has("ntfy") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# ntfy\n")
	txt.WriteString("ntfy_enabled: yes\n")

	return txt.String()
}

func (o *order) varsPostgresBackup() string {
	if o.has("borg") || o.hosting != "" {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# postgres::backups\n")
	txt.WriteString("postgres_backup_enabled: yes\n")
	txt.WriteString("postgres_backup_schedule: '@daily'\n")
	txt.WriteString("postgres_backup_keep_days: 7\n")
	txt.WriteString("postgres_backup_keep_weeks: 0\n")
	txt.WriteString("postgres_backup_keep_months: 0\n")

	return txt.String()
}

func (o *order) varsBorgBackup() string {
	if !o.has("borg") {
		return ""
	}
	var txt strings.Builder
	pub, priv := o.keygen()
	o.pass["borg"] = pub

	txt.WriteString("\n# borg\n")
	txt.WriteString("backup_borg_enabled: yes\n")
	txt.WriteString("backup_borg_location_repositories:\n")
	txt.WriteString("- \"" + o.get("borg-repository") + "\"\n")
	txt.WriteString("backup_borg_storage_encryption_passphrase: " + o.pwgen() + "\n")
	txt.WriteString("backup_borg_ssh_key_private: |\n")
	for line := range strings.SplitSeq(priv, "\n") {
		txt.WriteString("  " + line + "\n")
	}
	txt.WriteString("# " + pub + "\n")

	return txt.String()
}

func (o *order) varsEximRelay() string {
	var txt strings.Builder

	txt.WriteString("\n# exim-relay\n")
	if o.dkim["private"] != "" {
		txt.WriteString("exim_relay_dkim_privkey_contents: |\n")
		for line := range strings.SplitSeq(o.dkim["private"], "\n") {
			txt.WriteString("  " + line + "\n")
		}
	}

	if len(o.smtp) == 0 {
		return txt.String()
	}
	txt.WriteString("exim_relay_relay_use: yes\n")
	txt.WriteString("exim_relay_relay_auth: yes\n")
	txt.WriteString("exim_relay_relay_host_name: \"" + o.smtp["host"] + "\"\n")
	txt.WriteString("exim_relay_relay_host_port: \"" + o.smtp["port"] + "\"\n")
	txt.WriteString("exim_relay_relay_auth_username: " + o.smtp["login"] + "\n")
	txt.WriteString("exim_relay_relay_auth_password: \"" + o.smtp["password"] + "\"\n")
	txt.WriteString("exim_relay_sender_address: " + o.smtp["email"] + "\n")

	return txt.String()
}

func (o *order) varsMatrixRTC() string {
	if !o.has("matrix-rtc") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# matrix rtc\n")
	txt.WriteString("matrix_rtc_enabled: yes\n")

	return txt.String()
}

func (o *order) varsSynapse() string {
	var txt strings.Builder

	if o.has("sso") || o.has("synapse-sso") {
		txt.WriteString("\n# synapse::sso\n")
		txt.WriteString("matrix_synapse_oidc_enabled: yes\n")
		txt.WriteString("matrix_synapse_oidc_providers:\n")
		txt.WriteString(o.getOIDCConfig())
	}

	if o.has("synapse-workers") {
		txt.WriteString("\n# synapse::workers\n")
		txt.WriteString("matrix_synapse_workers_enabled: yes\n")
		txt.WriteString("matrix_synapse_workers_preset: specialized-workers\n")
	}

	if o.has("synapse-s3") || o.has("synapse-s3-storage") {
		txt.WriteString("\n# synapse::extensions::s3_storage_provider\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_enabled: yes\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_prefix: synapse\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_bucket: \"" + o.get("synapse-s3-bucket") + "\"\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_region_name: \"" + o.get("synapse-s3-region") + "\"\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_endpoint_url: \"" + o.get("synapse-s3-endpoint") + "\"\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_access_key_id: \"" + o.get("synapse-s3-access-key") + "\"\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_secret_access_key: \"" + o.get("synapse-s3-secret-key") + "\"\n")
		txt.WriteString("matrix_synapse_ext_synapse_s3_storage_provider_config_storage_class: STANDARD\n")
	}

	txt.WriteString("\n# synapse::extensions::shared_secret_auth\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_enabled: yes\n")
	txt.WriteString("matrix_synapse_ext_password_provider_shared_secret_auth_shared_secret: " + o.pwgen() + "\n")

	txt.WriteString(o.varsSynapseCredentials())

	return txt.String()
}

func (o *order) varsSynapseCredentials() string {
	var txt strings.Builder

	txt.WriteString("\n# synapse::credentials\n")
	txt.WriteString("matrix_synapse_macaroon_secret_key: " + o.pwgen() + "\n")
	txt.WriteString("matrix_synapse_password_config_pepper: " + o.pwgen() + "\n")
	txt.WriteString("matrix_coturn_turn_static_auth_secret: " + o.pwgen() + "\n")
	txt.WriteString("matrix_homeserver_generic_secret_key: \"{{ matrix_synapse_macaroon_secret_key }}\"\n")

	return txt.String()
}

func (o *order) varsSynapseAdmin() string {
	var txt strings.Builder
	txt.WriteString("\n# synapse-admin\n")
	txt.WriteString("matrix_synapse_admin_enabled: yes\n")
	txt.WriteString("matrix_synapse_admin_path_prefix: \"/admin\"\n")

	return txt.String()
}

func (o *order) varsUptimeKuma() string {
	if !o.has("uptime-kuma") {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("\n# uptime kuma\n")
	txt.WriteString("uptime_kuma_enabled: yes\n")
	txt.WriteString("uptime_kuma_hostname: \"kuma." + o.domain + "\"\n")

	return txt.String()
}

func (o *order) varsStaticFiles() string {
	if o.get("serve_base_domain") != "yes" {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("\n# matrix-static-files\n")
	txt.WriteString("matrix_static_files_container_labels_base_domain_enabled: " + o.get("serve_base_domain") + "\n")
	return txt.String()
}

func (o *order) varsRadicale() string {
	if !o.has("radicale") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# radicale\n")
	txt.WriteString("radicale_enabled: yes\n")
	txt.WriteString("radicale_hostname: \"radicale." + o.domain + "\"\n")

	return txt.String()
}

func (o *order) varsCinny() string {
	if !o.has("cinny") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# cinny\n")
	txt.WriteString("matrix_client_cinny_enabled: yes\n")

	return txt.String()
}

func (o *order) varsEtherpad() string {
	if !o.has("etherpad") {
		return ""
	}
	var txt strings.Builder

	o.login("etherpad admin", o.get("username"))
	txt.WriteString("\n# etherpad\n")
	txt.WriteString("etherpad_enabled: yes\n")
	txt.WriteString("etherpad_hostname: \"etherpad." + o.domain + "\"\n")
	txt.WriteString("etherpad_admin_username: \"" + o.get("username") + "\"\n")
	txt.WriteString("etherpad_admin_password: " + o.password("etherpad admin") + "\n")

	return txt.String()
}

func (o *order) varsFirezone() string {
	if !o.has("firezone") {
		return ""
	}
	var txt strings.Builder

	o.login("firezone", o.get("email"))
	txt.WriteString("\n# firezone\n")
	txt.WriteString("firezone_enabled: yes\n")
	txt.WriteString("firezone_hostname: \"firezone." + o.domain + "\"\n")
	txt.WriteString("firezone_default_admin_email: \"" + o.get("email") + "\"\n")
	txt.WriteString("firezone_default_admin_password: " + o.password("firezone") + "\n")
	txt.WriteString("firezone_database_encryption_key: \"" + o.base64bytesgen(32) + "\"\n")

	return txt.String()
}

func (o *order) varsFluffyChat() string {
	if !o.has("fluffychat") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# fluffychat\n")
	txt.WriteString("matrix_client_fluffychat_enabled: yes\n")

	return txt.String()
}

func (o *order) varsFunkwhale() string {
	if !o.has("funkwhale") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# funkwhale\n")
	txt.WriteString("funkwhale_enabled: yes\n")
	txt.WriteString("funkwhale_hostname: \"funkwhale." + o.domain + "\"\n")

	if o.has("funkwhale-s3-bucket") && o.has("funkwhale-s3-region") && o.has("funkwhale-s3-endpoint") && o.has("funkwhale-s3-access-key") && o.has("funkwhale-s3-secret-key") {
		txt.WriteString("funkwhale_aws_s3_region_name: \"" + o.get("funkwhale-s3-region") + "\"\n")
		txt.WriteString("funkwhale_aws_s3_endpoint_url: \"" + o.get("funkwhale-s3-endpoint") + "\"\n")
		txt.WriteString("funkwhale_aws_access_key_id: \"" + o.get("funkwhale-s3-access-key") + "\"\n")
		txt.WriteString("funkwhale_aws_secret_access_key: \"" + o.get("funkwhale-s3-secret-key") + "\"\n")
		txt.WriteString("funkwhale_aws_storage_bucket_name: \"" + o.get("funkwhale-s3-bucket") + "\"\n")
		txt.WriteString("funkwhale_aws_location: music\n")
	}

	return txt.String()
}

func (o *order) varsGoToSocial() string {
	if !o.has("gotosocial") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# gotosocial\n")
	txt.WriteString("gotosocial_enabled: yes\n")
	txt.WriteString("gotosocial_hostname: \"social." + o.domain + "\"\n")
	if o.has("gotosocial-s3-bucket") && o.has("gotosocial-s3-endpoint") && o.has("gotosocial-s3-access-key") && o.has("gotosocial-s3-secret-key") {
		txt.WriteString("gotosocial_environment_variables_extension: |\n")
		txt.WriteString("  GTS_STORAGE_BACKEND=s3\n")
		txt.WriteString("  GTS_STORAGE_S3_ENDPOINT=" + o.get("gotosocial-s3-endpoint") + "\n")
		txt.WriteString("  GTS_STORAGE_S3_BUCKET=" + o.get("gotosocial-s3-bucket") + "\n")
		txt.WriteString("  GTS_STORAGE_S3_ACCESS_KEY=" + o.get("gotosocial-s3-access-key") + "\n")
		txt.WriteString("  GTS_STORAGE_S3_SECRET_KEY=" + o.get("gotosocial-s3-secret-key") + "\n")
	}

	return txt.String()
}

func (o *order) varsElement() string {
	if !o.has("element-web") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# element\n")
	txt.WriteString("matrix_client_element_enabled: yes\n")

	return txt.String()
}

func (o *order) varsElementCall() string {
	if !o.has("element-call") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# element call\n")
	txt.WriteString("matrix_element_call_enabled: yes\n")

	return txt.String()
}

func (o *order) varsHydrogen() string {
	if !o.has("hydrogen") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# hydrogen\n")
	txt.WriteString("matrix_client_hydrogen_enabled: yes\n")

	return txt.String()
}

func (o *order) varsJitsi() string {
	if !o.has("jitsi") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# jitsi\n")
	txt.WriteString("jitsi_enabled: yes\n")
	txt.WriteString("# jitsi_enable_auth: yes\n")
	txt.WriteString("# jitsi_enable_guests: yes\n")
	txt.WriteString("jitsi_jvb_auth_password: " + o.pwgen() + "\n")
	txt.WriteString("jitsi_jibri_xmpp_password: " + o.pwgen() + "\n")
	txt.WriteString("jitsi_jibri_recorder_password: " + o.pwgen() + "\n")
	txt.WriteString("jitsi_jicofo_auth_password: " + o.pwgen() + "\n")
	txt.WriteString("# jitsi_prosody_auth_internal_accounts:\n")
	txt.WriteString("#  - username: \"" + o.get("username") + "\"\n")
	txt.WriteString("#    password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) varsLanguagetool() string {
	if !o.has("languagetool") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# languagetool\n")
	txt.WriteString("languagetool_enabled: yes\n")
	txt.WriteString("languagetool_hostname: \"languagetool." + o.domain + "\"\n")

	if o.has("languagetool-ngrams") {
		txt.WriteString("languagetool_ngrams_enabled: yes\n")
	}

	return txt.String()
}

func (o *order) varsLinkding() string {
	if !o.has("linkding") {
		return ""
	}
	var txt strings.Builder

	o.login("linkding", o.get("username"))
	txt.WriteString("\n# linkding\n")
	txt.WriteString("linkding_enabled: yes\n")
	txt.WriteString("linkding_hostname: \"linkding." + o.domain + "\"\n")
	txt.WriteString("linkding_superuser_username: \"" + o.get("username") + "\"\n")
	txt.WriteString("linkding_superuser_password: " + o.password("linkding") + "\n")

	return txt.String()
}

func (o *order) varsMiniflux() string {
	if !o.has("miniflux") {
		return ""
	}
	var txt strings.Builder

	o.login("miniflux", o.get("username"))
	txt.WriteString("\n# miniflux\n")
	txt.WriteString("miniflux_enabled: yes\n")
	txt.WriteString("miniflux_hostname: \"miniflux." + o.domain + "\"\n")
	txt.WriteString("miniflux_admin_login: \"" + o.get("username") + "\"\n")
	txt.WriteString("miniflux_admin_password: " + o.password("miniflux") + "\n")

	return txt.String()
}

func (o *order) varsPeertube() string {
	if !o.has("peertube") {
		return ""
	}
	var txt strings.Builder

	o.login("peertube", "root")
	txt.WriteString("\n# peertube\n")
	txt.WriteString("peertube_enabled: yes\n")
	txt.WriteString("peertube_hostname: \"peertube." + o.domain + "\"\n")
	txt.WriteString("peertube_config_secret: " + o.pwgen(64) + "\n")
	txt.WriteString("peertube_config_admin_email: \"" + o.get("email") + "\"\n")
	txt.WriteString("peertube_config_root_user_initial_password: " + o.password("peertube") + "\n")

	if o.has("peertube-s3-bucket") && o.has("peertube-s3-region") && o.has("peertube-s3-endpoint") && o.has("peertube-s3-access-key") && o.has("peertube-s3-secret-key") {
		txt.WriteString("peertube_config_object_storage_enabled: yes\n")
		txt.WriteString("peertube_config_object_storage_region: \"" + o.get("peertube-s3-region") + "\"\n")
		txt.WriteString("peertube_config_object_storage_endpoint: \"" + o.get("peertube-s3-endpoint") + "\"\n")
		txt.WriteString("peertube_config_object_storage_credentials_access_key_id: \"" + o.get("peertube-s3-access-key") + "\"\n")
		txt.WriteString("peertube_config_object_storage_credentials_secret_access_key: \"" + o.get("peertube-s3-secret-key") + "\"\n")
		txt.WriteString("peertube_config_object_storage_streaming_playlists_bucket_name: \"" + o.get("peertube-s3-bucket") + "\"\n")
		txt.WriteString("peertube_config_object_storage_streaming_playlists_prefix: playlists/\n")
		txt.WriteString("peertube_config_object_storage_web_videos_bucket_name: \"" + o.get("peertube-s3-bucket") + "\"\n")
		txt.WriteString("peertube_config_object_storage_web_videos_prefix: videos/\n")
	}

	return txt.String()
}

func (o *order) varsSchildiChat() string {
	if !o.has("schildichat") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# schildichat\n")
	txt.WriteString("matrix_client_schildichat_enabled: yes\n")

	return txt.String()
}

func (o *order) varsStats() string {
	if !o.has("stats") {
		return ""
	}
	var txt strings.Builder

	o.login("grafana", o.get("username"))
	txt.WriteString("\n# stats\n")
	txt.WriteString("grafana_enabled: yes\n")
	txt.WriteString("prometheus_enabled: yes\n")
	txt.WriteString("grafana_anonymous_access: no\n")
	txt.WriteString("prometheus_node_exporter_enabled: yes\n")
	txt.WriteString("prometheus_node_exporter_process_extra_arguments:\n")
	txt.WriteString("  - \"--collector.systemd\"\n")
	txt.WriteString("grafana_default_admin_user: \"" + o.get("username") + "\"\n")
	txt.WriteString("grafana_default_admin_password: " + o.password("grafana") + "\n")

	return txt.String()
}

func (o *order) varsVaultwarden() string {
	if !o.has("vaultwarden") {
		return ""
	}
	var txt strings.Builder
	o.login("vaultwarden admin", "_create yourself_")
	txt.WriteString("\n# vaultwarden\n")
	txt.WriteString("vaultwarden_enabled: yes\n")
	txt.WriteString("vaultwarden_hostname: \"vault." + o.domain + "\"\n")
	txt.WriteString("vaultwarden_config_admin_token: " + o.password("vaultwarden admin") + "\n")

	return txt.String()
}

func (o *order) varsBaibot() string {
	if !o.has("baibot") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::baibot\n")
	txt.WriteString("matrix_bot_baibot_enabled: yes\n")
	txt.WriteString("matrix_bot_baibot_config_user_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_baibot_config_user_encryption_recovery_passphrase: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_baibot_config_persistence_session_encryption_key: " + o.hexBytesGen(32) + "\n")
	txt.WriteString("matrix_bot_baibot_config_persistence_config_encryption_key: " + o.hexBytesGen(32) + "\n")

	return txt.String()
}

func (o *order) varsBuscarron() string {
	if !o.has("buscarron") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::buscarron\n")
	txt.WriteString("matrix_bot_buscarron_enabled: yes\n")
	txt.WriteString("matrix_bot_buscarron_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_buscarron_forms: []\n")

	return txt.String()
}

func (o *order) varsDraupnir() string {
	if !o.has("draupnir") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::draupnir\n")
	txt.WriteString("matrix_bot_draupnir_enabled: yes\n")
	txt.WriteString("matrix_bot_draupnir_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_draupnir_config_managementRoom: 'TBD'\n")

	return txt.String()
}

func (o *order) varsHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::honoroit\n")
	txt.WriteString("matrix_bot_honoroit_enabled: yes\n")
	txt.WriteString("matrix_bot_honoroit_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_bot_honoroit_roomid: 'TBD'\n")

	return txt.String()
}

func (o *order) varsMaubot() string {
	if !o.has("maubot") {
		return ""
	}
	var txt strings.Builder

	o.login("maubot admin", o.get("username"))
	txt.WriteString("\n# bots::maubot\n")
	txt.WriteString("matrix_bot_maubot_enabled: yes\n")
	txt.WriteString("matrix_bot_maubot_initial_password: " + o.password("@maubot:"+o.domain) + "\n")
	txt.WriteString("matrix_bot_maubot_admins:\n")
	txt.WriteString("  - \"" + o.get("username") + "\": " + o.password("maubot admin") + "\n")

	return txt.String()
}

func (o *order) varsReminder() string {
	if !o.has("reminder-bot") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bots::reminder\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_enabled: yes\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_reminders_timezone: \"" + o.get("reminder-bot-tz") + "\"\n")
	txt.WriteString("matrix_bot_matrix_reminder_bot_matrix_user_password: " + o.pwgen() + "\n")

	return txt.String()
}

func (o *order) varsBluesky() string {
	if !o.has("bluesky") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::bluesky\n")
	txt.WriteString("matrix_mautrix_bluesky_enabled: yes\n")

	return txt.String()
}

func (o *order) varsDiscord() string {
	if !o.has("discord") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::discord\n")
	txt.WriteString("matrix_mautrix_discord_enabled: yes\n")

	return txt.String()
}

func (o *order) varsEmail() string {
	if !o.has("postmoogle") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("\n# bridges::email\n")
	txt.WriteString("matrix_postmoogle_enabled: yes\n")
	txt.WriteString("matrix_postmoogle_password: " + o.pwgen() + "\n")
	txt.WriteString("matrix_postmoogle_data_secret: " + o.pwgen(32) + "\n")

	return txt.String()
}

func (o *order) varsSlack() string {
	if !o.has("slack") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::slack\n")
	txt.WriteString("matrix_mautrix_slack_enabled: yes\n")

	return txt.String()
}

func (o *order) varsSteam() string {
	if !o.has("steam") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::steam\n")
	txt.WriteString("matrix_steam_bridge_enabled: yes\n")

	return txt.String()
}

func (o *order) varsFacebook() string {
	if !o.has("facebook") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::facebook\n")
	txt.WriteString("matrix_mautrix_meta_messenger_enabled: yes\n")

	return txt.String()
}

func (o *order) varsGmessages() string {
	if !o.has("gmessages") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::gmessages\n")
	txt.WriteString("matrix_mautrix_gmessages_enabled: yes\n")

	return txt.String()
}

func (o *order) varsGooglechat() string {
	if !o.has("googlechat") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::googlechat\n")
	txt.WriteString("matrix_mautrix_googlechat_enabled: yes\n")

	return txt.String()
}

func (o *order) varsHeisenbridge() string {
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

func (o *order) varsInstagram() string {
	if !o.has("instagram") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::instagram\n")
	txt.WriteString("matrix_mautrix_meta_instagram_enabled: yes\n")

	return txt.String()
}

func (o *order) varsLinkedin() string {
	if !o.has("linkedin") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::linkedin\n")
	txt.WriteString("matrix_beeper_linkedin_enabled: yes\n")

	return txt.String()
}

func (o *order) varsSignal() string {
	if !o.has("signal") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::signal\n")
	txt.WriteString("matrix_mautrix_signal_enabled: yes\n")

	return txt.String()
}

func (o *order) varsTelegram() string {
	if !o.has("telegram") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::telegram\n")
	txt.WriteString("matrix_mautrix_telegram_enabled: yes\n")
	txt.WriteString("matrix_mautrix_telegram_api_id: \"" + o.get("telegram-api-id") + "\"\n")
	txt.WriteString("matrix_mautrix_telegram_api_hash: \"" + o.get("telegram-api-hash") + "\"\n")

	if o.has("telegram-bot-token") {
		txt.WriteString(`matrix_mautrix_telegram_bot_token: "` + o.get("telegram-bot-token") + `"` + "\n")
	}

	return txt.String()
}

func (o *order) varsTwitter() string {
	if !o.has("twitter") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::twitter\n")
	txt.WriteString("matrix_mautrix_twitter_enabled: yes\n")

	return txt.String()
}

func (o *order) varsWebhooks() string {
	if !o.has("webhooks") && !o.has("hookshot") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::hookshot\n")
	txt.WriteString("matrix_hookshot_enabled: yes\n")

	return txt.String()
}

func (o *order) varsWechat() string {
	if !o.has("wechat") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::wechat\n")
	txt.WriteString("matrix_wechat_enabled: yes\n")

	return txt.String()
}

func (o *order) varsWhatsapp() string {
	if !o.has("whatsapp") {
		return ""
	}
	var txt strings.Builder
	txt.WriteString("\n# bridges::whatsapp\n")
	txt.WriteString("matrix_mautrix_whatsapp_enabled: yes\n")

	return txt.String()
}

func (o *order) getOIDCConfig() string {
	id := strings.ToLower(o.get("sso-idp-id"))
	name := o.get("sso-idp-name")
	brand := strings.ToLower(o.get("sso-idp-brand"))
	issuer := o.get("sso-issuer")
	clientID := o.get("sso-client-id")
	clientSecret := o.get("sso-client-secret")

	provider := "default"
	keys := []string{id, name, brand}
	for _, key := range keys {
		if _, ok := oidcmap[key]; ok {
			provider = key
			break
		}
	}
	config := fmt.Sprintf(oidcmap[provider], id, name, brand, issuer, clientID, clientSecret)

	// special case: OIDC providers that require to use specific endpoints instead of auto-discovery
	issuerHost := "example.com"
	issuerURL, err := url.Parse(issuer)
	if err == nil {
		issuerHost = issuerURL.Host
		if provider == "microsoft" { // tenant id from issuer url: https://login.microsoftonline.com/<tenant id>/v2.0
			issuerHost = strings.Split(strings.TrimPrefix(issuerURL.Path, "/"), "/")[0]
		}
	}

	return strings.ReplaceAll(config, "example.com", issuerHost)
}
