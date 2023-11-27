package etkecc

import (
	"strings"
)

func (o *order) generateHosts() string {
	hasUser := o.has("ssh-user") && o.get("ssh-user") != "root"
	hasPass := o.has("ssh-password") && o.get("ssh-password") != ""
	hasPort := o.has("ssh-port") && o.get("ssh-port") != "22"

	var txt strings.Builder
	txt.WriteString(o.domain)
	txt.WriteString(" ")
	txt.WriteString("ansible_host=")
	txt.WriteString(o.get("ssh-host"))

	if hasUser {
		txt.WriteString(" ")
		txt.WriteString("ansible_user=")
		txt.WriteString(o.get("ssh-user"))
	}

	if hasPass {
		txt.WriteString(" ")
		txt.WriteString("ansible_become_password=")
		txt.WriteString(o.get("ssh-password"))
	}

	if hasPort {
		txt.WriteString(" ")
		txt.WriteString("ansible_port=")
		txt.WriteString(o.get("ssh-port"))
	}

	txt.WriteString(" ")
	txt.WriteString("ordered_at=")
	txt.WriteString(o.orderedAt.Format("2006-01-02_15:04:05"))

	return txt.String()
}
