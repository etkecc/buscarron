package etkecc

import (
	"encoding/json"
	"sort"
	"strings"
)

// hImage is "Ubuntu 22.04" image name
const hImage = "ubuntu-22.04"

var (
	hDomains = map[string]string{
		".etke.host":     "enTDpM8y67STAZcQMpmqr7",
		".kupo.email":    "X2dqkUMuGzqXenSx3Huz9T",
		".ma3x.chat":     "4Ys5JTRJ8Hyoip3UPDMVgj",
		".matrix.fan":    "x8RAPicWWX3kPurkL9ntVo",
		".matrix.town":   "Gn9RYjWvznHkLfzVdY6Gsa",
		".onmatrix.chat": "zVNMf3dur7oHP8dcGETZs",
	}
	hDNSPatch = map[string]string{
		"\n":            `\n`,
		"server IP":     "$SERVER_IP4",
		"ip4:SERVER_IP": "ip4:$SERVER_IP4",
	}
	// hLocations between human-readable name and actual name
	hLocations = map[string]string{
		"ashburn":     "ash",
		"falkenstein": "fsn1",
		"helsinki":    "hel1",
		"hillsboro":   "hil",
		"nuremberg":   "nbg1",
	}

	// hFirewall is the "matrix" firewall
	defaultFirewall = map[string]int{
		"firewall": 124003,
	}
	openFirewall = map[string]int{
		"firewall": 394512,
	}

	// hKeys is list of ssh keys names
	hKeys = []string{"first", "second", "third"}
)

type hVPSRequest struct {
	Name      string           `json:"name"`
	Size      string           `json:"server_type"`
	Image     string           `json:"image"`
	Firewalls []map[string]int `json:"firewalls"`
	SSHKeys   []string         `json:"ssh_keys"`
	Location  string           `json:"location"`
}

type hFirewallRequest struct {
	Name    string             `json:"name"`
	ApplyTo []hFirewallApplyTo `json:"apply_to"`
	Rules   []hFirewallRule    `json:"rules"`
}

type hFirewallApplyTo struct {
	Server hFirewallApplyToServer `json:"server"`
	Type   string                 `json:"type"`
}

type hFirewallApplyToServer struct {
	ID int `json:"id"`
}

type hFirewallRule struct {
	Description string   `json:"description"`
	Direction   string   `json:"direction"`
	Port        string   `json:"port"`
	Protocol    string   `json:"protocol"`
	SourceIPs   []string `json:"source_ips"`
}

type hDNSRecord struct {
	Subdomain string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	ZoneID    string `json:"zone_id"`
}

type hDNSRequest struct {
	Records    []hDNSRecord `json:"records"`
	WithMigadu bool         `json:"-"`
}

func (r *hDNSRequest) add(subdomain, rtype, value, zoneID string) *hDNSRequest {
	r.Records = append(r.Records, hDNSRecord{
		Subdomain: subdomain,
		Type:      rtype,
		Value:     value,
		ZoneID:    zoneID,
	})

	return r
}

func (o *order) generateHFirewallCommand() string {
	if !o.has("ssh-client-ips") || o.get("ssh-client-ips") == "N/A" {
		return ""
	}
	ips := []string{}
	for _, ip := range strings.Split(o.get("ssh-client-ips"), ",") {
		ips = append(ips, strings.TrimSpace(ip)+"/32")
	}
	req := &hFirewallRequest{
		Name: o.domain,
		ApplyTo: []hFirewallApplyTo{
			{
				Server: hFirewallApplyToServer{ID: 12345}, // special value to be replaced
				Type:   "server",
			},
		},
		Rules: []hFirewallRule{
			{
				Description: "SSH",
				Direction:   "in",
				Port:        "22",
				Protocol:    "tcp",
				SourceIPs:   ips,
			},
		},
	}
	return o.getHFirewallCurl(req)
}

func (o *order) generateHVPSCommand() string {
	location, ok := hLocations[strings.ToLower(o.get("turnkey-location"))]
	if !ok {
		location = "fsn1"
	}
	firewalls := []map[string]int{defaultFirewall}
	if o.get("ssh-client-ips") == "N/A" {
		firewalls = append(firewalls, openFirewall)
	}

	req := &hVPSRequest{
		Name:      o.domain,
		Size:      o.hosting,
		Image:     hImage,
		Firewalls: firewalls,
		SSHKeys:   hKeys,
		Location:  location,
	}

	return o.getHVPSCurl(req)
}

func (o *order) getHFirewallCurl(req *hFirewallRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(strings.ReplaceAll(string(reqb), "\"", "\\\""), "12345", "$SERVER_ID")

	var cmd strings.Builder
	cmd.WriteString(`curl -X "POST" "https://api.hetzner.cloud/v1/firewalls" `)
	cmd.WriteString(`-H "Content-Type: application/json" `)
	cmd.WriteString(`-H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" `)
	cmd.WriteString(`-d "`)
	cmd.WriteString(reqs)
	cmd.WriteString(`"`)
	cmd.WriteString("\n")
	return cmd.String()
}

func (o *order) getHVPSCurl(req *hVPSRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder
	cmd.WriteString("set -euxo pipefail\n")
	cmd.WriteString(`SERVER_INFO=$(curl -X "POST" "https://api.hetzner.cloud/v1/servers" `)
	cmd.WriteString(`-H "Content-Type: application/json" `)
	cmd.WriteString(`-H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" `)
	cmd.WriteString(`-d "`)
	cmd.WriteString(reqs)
	cmd.WriteString(`")`)
	cmd.WriteString("\n")

	cmd.WriteString(`SERVER_ID=$(echo $SERVER_INFO | jq -r '.server.id')`)
	cmd.WriteString("\n")
	cmd.WriteString(`SERVER_IP4=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.ip')`)
	cmd.WriteString("\n")
	cmd.WriteString(`SERVER_IP4_ID=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv4.id')`)
	cmd.WriteString("\n")
	cmd.WriteString(`SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')`)
	cmd.WriteString("\n")

	dnsPtrBody := strings.ReplaceAll(`{ "ip": "$SERVER_IP4", "dns_ptr": "matrix.`+o.domain+`" }`, "\"", "\\\"")
	cmd.WriteString(`curl -X "POST" "https://api.hetzner.cloud/v1/primary_ips/$SERVER_IP4_ID/actions/change_dns_ptr" `)
	cmd.WriteString(`-H "Content-Type: application/json" `)
	cmd.WriteString(`-H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD" `)
	cmd.WriteString(`-d "`)
	cmd.WriteString(dnsPtrBody)
	cmd.WriteString("\"\n")

	cmd.WriteString(`curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" `)
	cmd.WriteString(`-H "Content-Type: application/json" `)
	cmd.WriteString(`-H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"`)
	cmd.WriteString("\n")

	cmd.WriteString(o.generateHFirewallCommand())

	cmd.WriteString(`echo -e "---\n`)
	cmd.WriteString(`Hello,\n\n`)
	cmd.WriteString(`We've received your payment and have prepared a server for you. Its IP addresses are:\n\n`)
	cmd.WriteString(`- IPv4: $SERVER_IP4\n`)
	cmd.WriteString(`- IPv6: $SERVER_IP6\n`)
	dnsInstructions := o.adaptTurnkeyDNS()
	if dnsInstructions != "" {
		cmd.WriteString(dnsInstructions)
	}
	cmd.WriteString(`"`)
	cmd.WriteString("\n")
	return cmd.String()
}

func (o *order) adaptTurnkeyDNS() string {
	dnsEntries, internal := o.generateDNSInstructions()
	if internal {
		return ""
	}
	for from, to := range hDNSPatch {
		dnsEntries = strings.ReplaceAll(dnsEntries, from, to)
	}

	var msg strings.Builder
	msg.WriteString(dnsEntries)
	msg.WriteString(`\nPlease let us know when you're ready with the DNS`)
	if o.get("serve_base_domain") != "yes" {
		msg.WriteString(` and delegation redirects (you could just add the @ record pointing to the matrix server IP instead)`)
	}
	msg.WriteString(` configuration, so we can proceed with your server's setup.\n\n`)
	msg.WriteString(`Regards\n`)

	return msg.String()
}

//nolint:gocognit // TODO
func (o *order) generateHDNSCommand() string {
	req := &hDNSRequest{Records: []hDNSRecord{}}
	domain := o.domain
	subdomain := strings.Split(domain, ".")[0]
	suffix := "." + subdomain
	var zoneID string
	for sufix, zone := range hDomains {
		if strings.HasSuffix(domain, sufix) {
			zoneID = zone
			break
		}
	}

	serverIP := "$SERVER_IP4"
	serverIP6 := "$SERVER_IP6" // only for hosting
	if o.has("ssh-host") {
		serverIP = o.get("ssh-host")
	}

	req.
		add(subdomain, "A", serverIP, zoneID).
		add("matrix"+suffix, "A", serverIP, zoneID)

	if o.hosting != "" {
		req.add(subdomain, "AAAA", serverIP6, zoneID).
			add("matrix"+suffix, "AAAA", serverIP6, zoneID)
	}

	if o.has("service-email") {
		req.WithMigadu = true
		req.
			add(subdomain, "MX", "10 aspmx1.migadu.com.", zoneID).
			add(subdomain, "MX", "20 aspmx2.migadu.com.", zoneID).
			add("autoconfig"+suffix, "CNAME", "autoconfig.migadu.com.", zoneID).
			add("_autodiscover._tcp"+suffix, "SRV", "0 1 443 autodiscover.migadu.com", zoneID).
			add("key1._domainkey"+suffix, "CNAME", "key1."+domain+"._domainkey.migadu.com.", zoneID).
			add("key2._domainkey"+suffix, "CNAME", "key2."+domain+"._domainkey.migadu.com.", zoneID).
			add("key3._domainkey"+suffix, "CNAME", "key3."+domain+"._domainkey.migadu.com.", zoneID).
			add("_dmarc"+suffix, "TXT", "v=DMARC1; p=quarantine;", zoneID).
			add(subdomain, "TXT", "v=spf1 include:spf.migadu.com -all", zoneID).
			add(subdomain, "TXT", "hosted-email-verify=$MIGADU_VERIFICATION", zoneID)
	}

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)
	for _, key := range items {
		req.add(dnsmap[key]+suffix, "CNAME", "matrix."+o.domain+".", zoneID)
	}

	spf := "v=spf1 ip4:" + serverIP
	if o.hosting != "" {
		spf += " ip6:" + serverIP6
	}
	spf += " -all"

	// if there is no SMTP relay, we need to add SPF and DMARC records
	if len(o.smtp) == 0 {
		req.
			add(subdomain, "TXT", spf, zoneID).
			add("_dmarc"+suffix, "TXT", "v=DMARC1; p=quarantine;", zoneID)
	}

	// if there is email bridge, we need to add MX record and SPF/DMARC records (only if they were not added above)
	if o.has("email2matrix") || o.has("postmoogle") {
		req.add("matrix"+suffix, "MX", "0 matrix."+o.domain+".", zoneID)
		if len(o.smtp) > 0 {
			req.
				add("matrix"+suffix, "TXT", spf, zoneID).
				add("_dmarc.matrix"+suffix, "TXT", "v=DMARC1; p=quarantine;", zoneID)
		}
	}
	return o.getHDNSCurl(req)
}

func (o *order) getHDNSCurl(req *hDNSRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder

	if o.hosting == "" {
		cmd.WriteString("set -euxo pipefail\n")
	}
	if !o.has("ssh-host") && o.hosting == "" {
		cmd.WriteString("export SERVER_IP4=SERVER_IP\n")
	}
	if req.WithMigadu {
		cmd.WriteString("export MIGADU_VERIFICATION=CODE\n")
	}
	cmd.WriteString("curl -X \"POST\" \"https://dns.hetzner.com/api/v1/records/bulk\" ")
	cmd.WriteString("-H \"Content-Type: application/json\" ")
	cmd.WriteString("-H \"Auth-API-Token: $HETZNER_API_TOKEN\" ")
	cmd.WriteString("-d \"")
	cmd.WriteString(reqs)
	cmd.WriteString("\"\n")

	return cmd.String()
}
