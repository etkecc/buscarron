package etkecc

import (
	"sort"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/format"
)

func (o *order) generateOnboarding() {
	var txt strings.Builder

	txt.WriteString(o.generateOnboardingLinks())
	txt.WriteString(o.generateOnboardingBots())
	txt.WriteString(o.generateOnboardingBridges())

	txt.WriteString(o.generateOnboardingCredentials())
	txt.WriteString(o.generateOnboardingPayment())

	txt.WriteString(o.generateOnboardingAfter())

	content := format.RenderMarkdown(txt.String(), true, true)
	o.files = append(o.files,
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(content.Body),
			FileName:      "onboarding.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(content.Body)),
		},
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(content.FormattedBody),
			FileName:      "onboarding.html",
			ContentType:   "text/html",
			ContentLength: int64(len(content.FormattedBody)),
		},
	)
}

func (o *order) generateOnboardingLinks() string {
	var txt strings.Builder

	txt.WriteString("# " + o.t("links") + "\n\n")
	if !o.has("element-web") {
		txt.WriteString("* web: https://app.etke.cc\n")
	}
	txt.WriteString("* homeserver: https://matrix." + o.get("domain") + "\n")
	if o.has("synapse-admin") {
		txt.WriteString("* synapse-admin: https://matrix." + o.get("domain") + "/synapse-admin\n")
	}
	if o.has("etherpad") {
		url := "https://etherpad." + o.get("domain") + "/admin"
		if o.has("dimension") {
			url = "https://dimension." + o.get("domain") + "/etherpad/admin"
		}
		txt.WriteString("* etherpad admin: " + url + "\n")
	}

	items := []string{}
	for item := range dnsmap {
		if o.has(item) {
			if item == "etherpad" && o.has("dimension") {
				continue
			}
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, item := range items {
		txt.WriteString("* " + item + ": https://" + dnsmap[item] + "." + o.get("domain") + "\n")
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
	txt.WriteString("# " + o.t("bots") + "\n\n")
	items := []string{}
	for item := range botmap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bot := range items {
		if o.has(bot) {
			txt.WriteString("* " + bot + ": " + botmap[bot] + ":" + o.get("domain") + "\n")
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
	txt.WriteString("# " + o.t("bridges") + "\n\n")
	items := []string{}
	for item := range bridgemap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bridge := range items {
		txt.WriteString("* " + bridge + ": " + bridgemap[bridge] + ":" + o.get("domain") + "\n")
	}
	txt.WriteString("\n\n")
	txt.WriteString("> https://etke.cc/" + o.get("lang") + "/bridges - " + o.t("auth_instructions") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingCredentials() string {
	var txt strings.Builder

	// hacky way to simplify next loop
	mxpass := o.pass["matrix"]
	delete(o.pass, "matrix")

	txt.WriteString("# " + o.t("credentials") + "\n\n")
	txt.WriteString("* mxid: @" + o.get("username") + ":" + o.get("domain") + "\n")
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
	txt.WriteString("> " + o.t("in_case_of_issues") + ": https://etke.cc/" + o.get("lang") + "/help\n\n")
	o.pass["matrix"] = mxpass

	return txt.String()
}

func (o *order) generateOnboardingPayment() string {
	if o.get("type") == "turnkey" {
		return ""
	}
	membership := "Maintenance"
	if o.has("service-email") {
		membership = "Maintenance+Email"
	}

	var txt strings.Builder
	txt.WriteString("# " + o.t("payment") + "\n\n")
	txt.WriteString(o.t("buy_setup") + "(https://etke.cc/setup)")
	txt.WriteString(" " + o.t("and") + " " + o.t("join_the") + " **" + membership + "** " + o.t("membership_on") + " [https://etke.cc/membership](https://etke.cc/membership).\n")
	txt.WriteString("\n" + o.t("note_pwyw") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfter() string {
	has := (o.has("etherpad") && o.has("dimension")) || o.has("honoroit")
	if !has {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("# " + o.t("steps_after_setup") + "\n\n")
	txt.WriteString(o.generateOnboardingAfterEtherpad())
	txt.WriteString(o.generateOnboardingAfterBuscarron())
	txt.WriteString(o.generateOnboardingAfterHonoroit())

	return txt.String()
}

func (o *order) generateOnboardingAfterEtherpad() string {
	if !(o.has("etherpad") && o.has("dimension")) {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### etherpad\n\n")
	txt.WriteString("1. " + o.t("as_etherpad_1") + "\n")
	txt.WriteString("2. " + o.t("as_etherpad_2") + "\n")
	txt.WriteString("3. " + o.t("as_etherpad_3") + " `dimension." + o.get("domain") + "`\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterBuscarron() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### buscarron\n\n")
	txt.WriteString("1. " + o.t("as_buscarron_1") + "\n")
	txt.WriteString("2. " + o.t("as_buscarron_2") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### honoroit\n\n")
	txt.WriteString("1. " + o.t("as_honoroit_1") + "\n")
	txt.WriteString("2. " + o.t("as_honoroit_2") + "\n\n")

	txt.WriteString(o.t("as_honoroit_3") + "\n\n")

	return txt.String()
}
