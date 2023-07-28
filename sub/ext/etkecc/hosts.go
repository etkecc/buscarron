package etkecc

import "strings"

func (o *order) generateHosts() string {
	if !o.has("ssh-host") {
		return ""
	}
	user := o.get("ssh-user")
	pass := o.get("ssh-password")
	port := o.get("ssh-port")

	var txt strings.Builder
	txt.WriteString(o.get("domain"))
	txt.WriteString(" ")
	txt.WriteString("ansible_host=")
	txt.WriteString(o.get("ssh-host"))

	if user != "" && user != "root" {
		txt.WriteString(" ")
		txt.WriteString("ansible_user=")
		txt.WriteString(user)
	}

	if pass != "" {
		txt.WriteString(" ")
		txt.WriteString("ansible_become_password=")
		txt.WriteString(pass)
	}

	if port != "" && port != "22" {
		txt.WriteString(" ")
		txt.WriteString("ansible_port=")
		txt.WriteString(port)
	}

	txt.WriteString("\n")
	return txt.String()
}
