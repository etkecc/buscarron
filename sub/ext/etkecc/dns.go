package etkecc

import (
	"encoding/json"
	"sort"
	"strings"
)

type bulkRecordsRequest struct {
	Records []dnsRecord `json:"records"`
}

func (r *bulkRecordsRequest) add(subdomain, rtype, value string) *bulkRecordsRequest {
	r.Records = append(r.Records, dnsRecord{
		Subdomain: subdomain,
		Type:      rtype,
		Value:     value,
		ZoneID:    "$HETZNER_ZONE_ID",
	})

	return r
}

type dnsRecord struct {
	Subdomain string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	ZoneID    string `json:"zone_id"`
}

func (o *order) generateDNSInstructions() string {
	if o.get("domain-type") == "subdomain" {
		return o.generateDNSCommand()
	}
	dns := "\n" + o.t("dns_add_entries") + ":\n"
	if o.get("serve_base_domain") == "yes" {
		dns += strings.Join([]string{"@", "A record", "server IP\n"}, "\t")
	}
	dns += strings.Join([]string{"matrix", "A record", "server IP\n"}, "\t")

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)

	for _, key := range items {
		dns += strings.Join([]string{dnsmap[key], "CNAME record", "matrix." + o.get("domain") + "\n"}, "\t")
	}

	if o.has("email2matrix") || o.has("postmoogle") {
		dns += strings.Join([]string{"matrix", "MX record", "matrix." + o.get("domain") + "\n"}, "\t")
		dns += strings.Join([]string{"matrix", "TXT record", "v=spf1 ip4:SERVER_IP -all\n"}, "\t")
		dns += strings.Join([]string{"_dmarc.matrix", "TXT record", "v=DMARC1; p=quarantine;\n"}, "\t")
	}

	return dns
}

func (o *order) generateDNSCommand() string {
	req := &bulkRecordsRequest{Records: []dnsRecord{}}
	domain := o.get("domain")
	subdomain := strings.Split(domain, ".")[0]
	suffix := "." + subdomain

	req.
		add(subdomain, "A", "$HETZNER_SERVER_IP").
		add("matrix"+suffix, "A", "$HETZNER_SERVER_IP").
		add(subdomain, "MX", "10 aspmx1.migadu.com.").
		add(subdomain, "MX", "20 aspmx2.migadu.com.").
		add("autoconfig"+suffix, "CNAME", "autoconfig.migadu.com.").
		add("_autodiscover._tcp"+suffix, "SRV", "0 1 443 autodiscover.migadu.com").
		add("key1._domainkey"+suffix, "CNAME", "key1."+domain+"._domainkey.migadu.com.").
		add("key2._domainkey"+suffix, "CNAME", "key2."+domain+"._domainkey.migadu.com.").
		add("key3._domainkey"+suffix, "CNAME", "key3."+domain+"._domainkey.migadu.com.").
		add("_dmarc"+suffix, "TXT", "v=DMARC1; p=quarantine;").
		add(subdomain, "TXT", "v=spf1 include:spf.migadu.com -all")

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)
	for _, key := range items {
		req.add(dnsmap[key]+suffix, "CNAME", "matrix."+o.get("domain"))
	}

	if o.has("email2matrix") || o.has("postmoogle") {
		req.
			add("matrix"+suffix, "MX", "0 matrix."+o.get("domain")).
			add("matrix"+suffix, "TXT", "v=spf1 ip4:$HETZNER_SERVER_IP -all").
			add("_dmarc.matrix"+suffix, "TXT", "v=DMARC1; p=quarantine;")
	}
	return o.getDNSCurl(req)
}

func (o *order) getDNSCurl(req *bulkRecordsRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder
	cmd.WriteString("export HETZNER_SERVER_IP=SERVER_IP\n")
	cmd.WriteString("curl -X \"POST\" \"https://dns.hetzner.com/api/v1/records/bulk\" ")
	cmd.WriteString("-H \"Content-Type: application/json\" ")
	cmd.WriteString("-H \"Auth-API-Token: $HETZNER_API_TOKEN\" ")
	cmd.WriteString("-d \"")
	cmd.WriteString(reqs)
	cmd.WriteString("\"\n")

	return cmd.String()
}
