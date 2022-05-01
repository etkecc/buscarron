package etkecc

import (
	"sort"
	"strings"
)

func (o *order) generateDNS() string {
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

	return dns
}
