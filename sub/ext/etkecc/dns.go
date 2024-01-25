package etkecc

import (
	"sort"
	"strings"
)

//nolint:gocognit // TODO
func (o *order) generateDNSInstructions() (string, bool) {
	if o.subdomain {
		return o.generateHDNSCommand(), true
	}

	serverIP := "server IP"
	if o.has("ssh-host") {
		serverIP = o.get("ssh-host")
	}

	dns := "\nPlease, add the following DNS entries"
	if o.v.NS(o.domain, "cloudflare.com") {
		dns += " (ensure that the CloudFlare proxy is disabled, as it's known to cause issues with Matrix Federation)"
	}
	dns += ":\n\n"
	if o.get("serve_base_domain") == "yes" {
		dns += strings.Join([]string{"- @", "A record", serverIP + "\n"}, "\t")
		if o.hosting != "" {
			dns += strings.Join([]string{"- @", "AAAA record", "$SERVER_IP6\n"}, "\t")
		}
	}
	dns += strings.Join([]string{"- matrix", "A record", serverIP + "\n"}, "\t")
	if o.hosting != "" {
		dns += strings.Join([]string{"- matrix", "AAAA record", "$SERVER_IP6\n"}, "\t")
	}

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)

	for _, key := range items {
		dns += strings.Join([]string{"- " + dnsmap[key], "CNAME record", "matrix." + o.domain + "\n"}, "\t")
	}

	spf := "v=spf1 ip4:" + serverIP
	if o.hosting != "" {
		spf += " ip6:$SERVER_IP6"
	}
	spf += " -all\n"

	// if there is no SMTP relay, we need to add SPF and DMARC records
	if len(o.smtp) == 0 {
		dns += strings.Join([]string{"- matrix", "TXT record", spf}, "\t")
		dns += strings.Join([]string{"- _dmarc.matrix", "TXT record", "v=DMARC1; p=quarantine;\n"}, "\t")
	}

	// if there is email bridge, we need to add MX record and SPF/DMARC records (only if they were not added above)
	if o.has("email2matrix") || o.has("postmoogle") {
		dns += strings.Join([]string{"- matrix", "MX record", "matrix." + o.domain + "\n"}, "\t")
		if len(o.smtp) > 0 {
			dns += strings.Join([]string{"- matrix", "TXT record", spf}, "\t")
			dns += strings.Join([]string{"- _dmarc.matrix", "TXT record", "v=DMARC1; p=quarantine;\n"}, "\t")
		}
	}

	return dns, false
}
