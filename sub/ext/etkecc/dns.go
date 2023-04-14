package etkecc

import (
	"sort"
	"strings"
)

func (o *order) generateDNSInstructions() (string, bool) {
	if o.get("domain-type") == "subdomain" {
		return o.generateHDNSCommand(), true
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

	return dns, false
}
