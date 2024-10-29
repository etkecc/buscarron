package etkecc

import (
	"context"
	"sort"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
)

func (o *order) generateDNSInstructions(ctx context.Context) string {
	log := o.logger(ctx)
	log.Info().Msg("generating DNS instructions")
	if o.subdomain || o.hosting != "" {
		log.Info().Msg("skipping DNS instructions for subdomain or hosting")
		return ""
	}

	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateDNSInstructions")
	defer span.Finish()

	serverIP := "server IP"
	if o.has("ssh-host") {
		serverIP = o.get("ssh-host")
	}

	dns := "\nPlease, add the following DNS entries"
	if o.v.NS(o.domain, "cloudflare.com") {
		dns += " (ensure that the CloudFlare proxy is disabled, as it's known to cause issues with Matrix Federation)"
	}
	dns += ":\n\n"

	for _, record := range o.generateDNSRecords("@", "", serverIP, "") {
		record = strings.Trim(record, `"`)
		parts := strings.Split(record, ",")
		parts[1] += " record"
		dns += "- " + strings.Join(parts, "\t") + "\n"
	}

	log.Info().Msg("DNS instructions have been generated")
	return dns
}

func (o *order) generateDNSRecords(domainRecord, suffix, serverIPv4, serverIPv6 string) []string {
	records := []string{}
	if o.has("serve_base_domain") || o.subdomain {
		records = append(records, domainRecord+",A,"+serverIPv4)
		if serverIPv6 != "" {
			records = append(records, domainRecord+",AAAA,"+serverIPv6)
		}
	}
	records = append(records, "matrix"+suffix+",A,"+serverIPv4)
	if serverIPv6 != "" {
		records = append(records, "matrix"+suffix+",AAAA,"+serverIPv6)
	}

	items := []string{}
	for key := range dnsmap {
		if o.has(key) {
			items = append(items, key)
		}
	}
	sort.Strings(items)
	for _, key := range items {
		records = append(records, dnsmap[key]+suffix+",CNAME,matrix."+o.domain+".")
	}

	spf := o.generateDNSSPF(serverIPv4, serverIPv6)
	records = append(records,
		"matrix"+suffix+",TXT,"+spf,
		"_dmarc.matrix"+suffix+",TXT,v=DMARC1; p=quarantine;",
	)
	if o.dkim["record"] != "" {
		records = append(records, "default._domainkey.matrix"+suffix+",TXT,"+o.dkim["record"])
	}
	if o.has("postmoogle") {
		records = append(records, "matrix"+suffix+",MX,0 matrix."+o.domain+".")
	}

	if o.has("service-email") {
		records = append(records,
			domainRecord+",MX,10 aspmx1.migadu.com.",
			domainRecord+",MX,20 aspmx2.migadu.com.",
			domainRecord+",TXT,v=spf1 include:spf.migadu.com -all",
			"autoconfig"+suffix+",CNAME,autoconfig.migadu.com.",
			"key1._domainkey"+suffix+",CNAME,key1."+o.domain+"._domainkey.migadu.com.",
			"key2._domainkey"+suffix+",CNAME,key2."+o.domain+"._domainkey.migadu.com.",
			"key3._domainkey"+suffix+",CNAME,key3."+o.domain+"._domainkey.migadu.com.",
			"_dmarc"+suffix+",TXT,v=DMARC1; p=quarantine;",
			"_autodiscover._tcp"+suffix+",SRV,0 1 443 autodiscover.migadu.com",
		)
	}

	for i, record := range records {
		records[i] = `"` + record + `"`
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
