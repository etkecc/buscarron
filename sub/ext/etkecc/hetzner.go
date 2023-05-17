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
		".etke.host": "enTDpM8y67STAZcQMpmqr7",
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
	hFirewall = map[string]int{
		"firewall": 124003,
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

func (o *order) generateHVPSCommand() string {
	location, ok := hLocations[strings.ToLower(o.get("turnkey-location"))]
	if !ok {
		location = "fsn1"
	}
	var size string
	sizeParts := strings.Split(o.get("turnkey"), "-")
	if len(sizeParts) < 2 {
		size = "cx11"
	} else {
		size = sizeParts[1]
	}

	req := &hVPSRequest{
		Name:      o.get("domain"),
		Size:      size,
		Image:     hImage,
		Firewalls: []map[string]int{hFirewall},
		SSHKeys:   hKeys,
		Location:  location,
	}

	return o.getHVPSCurl(req)
}

func (o *order) getHVPSCurl(req *hVPSRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder
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
	cmd.WriteString(`SERVER_IP6=$(echo $SERVER_INFO | jq -r '.server.public_net.ipv6.ip' | sed -e 's|/64|1|g')`)
	cmd.WriteString("\n")

	cmd.WriteString(`curl -X "POST" "https://api.hetzner.cloud/v1/servers/$SERVER_ID/actions/enable_backup" `)
	cmd.WriteString(`-H "Content-Type: application/json" `)
	cmd.WriteString(`-H "Authorization: Bearer $HETZNER_API_TOKEN_CLOUD"`)
	cmd.WriteString("\n")

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
	msg.WriteString(`\nIf you care about IPv6, feel free to configure additional AAAA records in the steps mentioning A records above.\n\n`)
	msg.WriteString(`Let us know when you're ready with the DNS configuration, so we can proceed with your server's setup.\n\n`)
	msg.WriteString(`Regards\n`)

	return msg.String()
}

func (o *order) generateHDNSCommand() string {
	req := &hDNSRequest{Records: []hDNSRecord{}}
	domain := o.get("domain")
	subdomain := strings.Split(domain, ".")[0]
	suffix := "." + subdomain
	var zoneID string
	for sufix, zone := range hDomains {
		if strings.HasSuffix(domain, sufix) {
			zoneID = zone
			break
		}
	}

	req.
		add(subdomain, "A", "$HETZNER_SERVER_IP", zoneID).
		add("matrix"+suffix, "A", "$HETZNER_SERVER_IP", zoneID)

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
		req.add(dnsmap[key]+suffix, "CNAME", "matrix."+o.get("domain")+".", zoneID)
	}

	if o.has("email2matrix") || o.has("postmoogle") {
		req.
			add("matrix"+suffix, "MX", "0 matrix."+o.get("domain")+".", zoneID).
			add("matrix"+suffix, "TXT", "v=spf1 ip4:$HETZNER_SERVER_IP -all", zoneID).
			add("_dmarc.matrix"+suffix, "TXT", "v=DMARC1; p=quarantine;", zoneID)
	}
	return o.getHDNSCurl(req)
}

func (o *order) getHDNSCurl(req *hDNSRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder
	cmd.WriteString("export HETZNER_SERVER_IP=SERVER_IP\n")
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
