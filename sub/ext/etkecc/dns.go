package etkecc

import (
	"context"
	"sort"
	"strings"

	"gitlab.com/etke.cc/buscarron/utils"
)

func (o *order) generateDNSInstructions(ctx context.Context) (string, bool) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateDNSInstructions")
	defer span.Finish()

	if o.subdomain {
		return o.generateHDNSCommand(span.Context()), true
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

	for _, record := range o.generateDNSRecords(serverIP) {
		dns += "- " + strings.Join(record, "\t") + "\n"
	}

	return dns, false
}

func (o *order) generateDNSRecords(serverIP string) [][]string {
	var serverIPv6 string
	if o.hosting != "" {
		serverIPv6 = "$SERVER_IP6"
	}

	records := [][]string{}
	if o.get("serve_base_domain") == "yes" {
		records = append(records, []string{"@", "A record", serverIP})
		if o.hosting != "" {
			records = append(records, []string{"@", "AAAA record", "$SERVER_IP6"})
		}
	}
	records = append(records, []string{"matrix", "A record", serverIP})
	if o.hosting != "" {
		records = append(records, []string{"matrix", "AAAA record", "$SERVER_IP6"})
	}

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)

	for _, key := range items {
		records = append(records, []string{dnsmap[key], "CNAME record", "matrix." + o.domain})
	}

	spf := o.generateDNSSPF(serverIP, serverIPv6)
	// if there is no SMTP relay, we need to add SPF and DMARC records
	if len(o.smtp) == 0 {
		records = append(records,
			[]string{"matrix", "TXT record", spf},
			[]string{"_dmarc.matrix", "TXT record", "v=DMARC1; p=quarantine;"},
		)
	}

	// if there is email bridge, we need to add MX record and SPF/DMARC records (only if they were not added above)
	if o.has("email2matrix") || o.has("postmoogle") {
		records = append(records, []string{"matrix", "MX record", "matrix." + o.domain})
		if len(o.smtp) > 0 {
			records = append(records,
				[]string{"matrix", "TXT record", spf},
				[]string{"_dmarc.matrix", "TXT record", "v=DMARC1; p=quarantine;"},
			)
		}
	}

	if o.has("service-email") {
		records = append(records,
			[]string{"@", "MX record", "10 aspmx1.migadu.com"},
			[]string{"@", "MX record", "20 aspmx2.migadu.com"},
			[]string{"@", "TXT record", "v=spf1 include:spf.migadu.com -all"},
			[]string{"autoconfig", "CNAME record", "autoconfig.migadu.com"},
			[]string{"key1._domainkey", "CNAME record", "key1." + o.domain + "._domainkey.migadu.com"},
			[]string{"key2._domainkey", "CNAME record", "key2." + o.domain + "._domainkey.migadu.com"},
			[]string{"key3._domainkey", "CNAME record", "key3." + o.domain + "._domainkey.migadu.com"},
			[]string{"_dmarc", "TXT record", "v=DMARC1; p=quarantine;"},
			[]string{"_autodiscover._tcp", "SRV record", "0 1 443 autodiscover.migadu.com"},
		)
	}

	return records
}

func (o *order) generateDNSSPF(serverIPv4 string, serverIPv6 ...string) string {
	spf := "v=spf1 ip4:" + serverIPv4
	if len(serverIPv6) > 0 && serverIPv6[0] != "" {
		spf += " ip6:" + serverIPv6[0]
	}
	spf += " -all"

	return spf
}
