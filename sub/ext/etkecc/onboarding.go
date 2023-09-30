package etkecc

import (
	"sort"
	"strings"

	"maunium.net/go/mautrix"
)

func (o *order) generateOnboarding() {
	var txt strings.Builder

	txt.WriteString(o.generateOnboardingIntro())
	txt.WriteString(o.generateOnboardingLinks())
	txt.WriteString(o.generateOnboardingBots())
	txt.WriteString(o.generateOnboardingBridges())

	txt.WriteString(o.generateOnboardingCredentials())

	txt.WriteString(o.generateOnboardingAfter())

	txt.WriteString(o.generateOnboardingOutro())

	text := txt.String()
	o.files = append(o.files,
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(text),
			FileName:      "onboarding.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(text)),
		},
	)
}

func (o *order) generateOnboardingIntro() string {
	var txt strings.Builder
	txt.WriteString("Hello!\n\n")
	txt.WriteString("Your Matrix server is ready, here is your onboarding list:\n\n")

	return txt.String()
}

func (o *order) generateOnboardingLinks() string {
	var txt strings.Builder

	txt.WriteString("# Links\n\n")
	if !o.has("element-web") {
		txt.WriteString("* web: " + link("app.etke.cc") + "\n")
	}
	txt.WriteString("* homeserver: " + link("matrix."+o.domain) + "\n")
	if o.has("synapse-admin") {
		txt.WriteString("* synapse-admin: " + link("matrix."+o.domain+"/synapse-admin") + "\n")
	}
	if o.has("etherpad") {
		txt.WriteString("* etherpad admin: " + link("etherpad."+o.domain+"/admin") + "\n")
	}
	if o.has("vaultwarden") {
		txt.WriteString("* vaultwarden admin:" + link("vault."+o.domain+"/admin") + "\n")
	}

	items := []string{}
	for item := range dnsmap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, item := range items {
		txt.WriteString("* " + item + ": " + link(dnsmap[item]+"."+o.domain) + "\n")
	}
	txt.WriteString("\n\n")

	return txt.String()
}

func (o *order) generateOnboardingBots() string {
	skip := true
	for bot := range botmap {
		if o.has(bot) {
			skip = false
			break
		}
	}
	if skip {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("# Bots\n\n")
	items := []string{}
	for item := range botmap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bot := range items {
		if o.has(bot) {
			txt.WriteString("* " + bot + ": " + matrixLink(botmap[bot]+":"+o.domain) + "\n")
		}
	}
	txt.WriteString("\n\n")

	return txt.String()
}

func (o *order) generateOnboardingBridges() string {
	skip := true
	for bridge := range bridgemap {
		if o.has(bridge) {
			skip = false
			break
		}
	}
	if skip {
		return ""
	}

	var txt strings.Builder
	txt.WriteString("# Bridges\n\n")
	items := []string{}
	for item := range bridgemap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bridge := range items {
		txt.WriteString("* " + bridge + ": " + matrixLink(bridgemap[bridge]+":"+o.domain) + "\n")
	}
	txt.WriteString("\n\n")
	txt.WriteString("> auth instruction and support: " + link("etke.cc/help/bridges") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingCredentials() string {
	var txt strings.Builder

	// hacky way to simplify next loop
	mxpass := o.pass["matrix"]
	delete(o.pass, "matrix")

	txt.WriteString("# Credentials\n\n")
	txt.WriteString("* mxid: " + matrixLink("@"+o.get("username")+":"+o.domain) + "\n")
	txt.WriteString("* username: " + o.get("username") + "\n")
	txt.WriteString("* password: " + mxpass + "\n")
	items := []string{}
	for item := range o.pass {
		items = append(items, item)
	}
	sort.Strings(items)
	for _, name := range items {
		txt.WriteString("* " + name + " password: " + o.pass[name] + "\n")
	}
	txt.WriteString("\n\n")
	txt.WriteString("> in case of any issues: " + link("etke.cc/help") + "\n\n")
	o.pass["matrix"] = mxpass

	return txt.String()
}

func (o *order) generateOnboardingOutro() string {
	var txt strings.Builder
	txt.WriteString("Happy Matrixing!\n\n")

	txt.WriteString("PS: To help with your first steps in the world of Matrix, we've built the " + link("MatrixRooms.info") + " search engine. You can use it to discover rooms over the Matrix Federation and find amazing communities you wish to join!\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfter() string {
	if !o.has("honoroit") || !o.has("buscarron") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("# Steps after the setup\n\n")
	txt.WriteString(o.generateOnboardingAfterBuscarron())
	txt.WriteString(o.generateOnboardingAfterHonoroit())

	return txt.String()
}

func (o *order) generateOnboardingAfterBuscarron() string {
	if !o.has("buscarron") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### buscarron\n\n")
	txt.WriteString("1. Create (encrypted) Matrix room(-s) for the forms you want to have and invite the buscarron user to all of them\n")
	txt.WriteString("2. Send the list of (form name, Matrix room id, redirect URL after submission) to " + matrixLink("@support:etke.cc") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### honoroit\n\n")
	txt.WriteString("1. Create a matrix room (encryption supported) and invite the honoroit user into it\n")
	txt.WriteString("2. Send the room id to " + matrixLink("@support:etke.cc") + "\n\n")

	return txt.String()
}

func matrixLink(id string) string {
	return "[" + id + "](https://matrix.to/#/" + id + ")"
}

func link(address string) string {
	return "[" + address + "](https://" + address + ")"
}
