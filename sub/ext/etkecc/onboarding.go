package etkecc

import (
	"sort"
	"strings"

	"github.com/russross/blackfriday/v2"
	"maunium.net/go/mautrix"
)

func (o *order) generateOnboarding() {
	var txt strings.Builder

	txt.WriteString(o.generateOnboardingLinks())
	txt.WriteString(o.generateOnboardingBots())
	txt.WriteString(o.generateOnboardingBridges())

	txt.WriteString(o.generateOnboardingCredentials())
	txt.WriteString(o.generateOnboardingPayment())

	txt.WriteString(o.generateOnboardingAfter())

	htmlBytes := blackfriday.Run([]byte(txt.String()), bfExtsOpt, bfRendererOpt)
	html := strings.TrimRight(string(htmlBytes), "\n")
	html = htmlPRegex.ReplaceAllString(html, "$1")

	o.files = append(o.files,
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(txt.String()),
			FileName:      "onboarding.md",
			ContentType:   "text/markdown",
			ContentLength: int64(txt.Len()),
		},
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(html),
			FileName:      "onboarding.html",
			ContentType:   "text/html",
			ContentLength: int64(len(html)),
		},
	)
}

func (o *order) generateOnboardingLinks() string {
	var txt strings.Builder

	txt.WriteString("# links\n\n")
	txt.WriteString("* homeserver: https://matrix." + o.get("domain") + "\n")
	txt.WriteString("* synapse-admin: https://matrix." + o.get("domain") + "/synapse-admin\n")
	if o.has("etherpad") {
		txt.WriteString("* etherpad admin: https://dimension." + o.get("domain") + "/etherpad/admin\n")
	}
	if o.has("softserve") {
		txt.WriteString("* ssh git: ssh matrix." + o.get("domain") + ":23231\n")
	}
	items := []string{}
	for item := range dnsmap {
		if o.has(item) {
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
	txt.WriteString("# bots\n\n")
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
	txt.WriteString("# bridges\n\n")
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
	txt.WriteString("> https://etke.cc/bridges - auth instructions\n\n")

	return txt.String()
}

func (o *order) generateOnboardingCredentials() string {
	var txt strings.Builder

	// hacky way to simplify next loop
	mxpass := o.pass["matrix"]
	delete(o.pass, "matrix")

	txt.WriteString("# credentials\n\n")
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
	txt.WriteString("> in case of any issues: @support:etke.cc\n\n")
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
	txt.WriteString("# payment\n\n")
	txt.WriteString("Please, [buy the Setup item](https://etke.cc/setup)")
	if o.has("service-maintenance") || o.has("service-email") {
		txt.WriteString(" and join the **" + membership + "** membership on [https://etke.cc/membership](https://etke.cc/membership).\n")
	} else {
		txt.WriteString(".\nIf you want an ongoing maintenance of your server (host/system maintenance, matrix components maintenance, updates and reconfiguration), join the **Maintenance** membership on [https://etke.cc/membership](https://etke.cc/membership).\nPlease, clarify now if you want to join the maintenance membership, because if not - we will remove any your configuration and credentials from our side.\n")
	}
	txt.WriteString("\n> **NOTE**: all prices are based on [Pay What You Want](https://en.wikipedia.org/wiki/Pay_what_you_want) model.\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfter() string {
	has := o.has("email2matrix") || o.has("etherpad") || o.has("honoroit")
	if !has {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("# steps after the setup\n\n")
	txt.WriteString(o.generateOnboardingAfterEmail2Matrix())
	txt.WriteString(o.generateOnboardingAfterEtherpad())
	txt.WriteString(o.generateOnboardingAfterHonoroit())

	return txt.String()
}

func (o *order) generateOnboardingAfterEmail2Matrix() string {
	if !o.has("email2matrix") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### email2matrix\n\n")
	txt.WriteString("1. Add new MX record on matrix subdomain that will point... to matrix." + o.get("domain") + " (looks odd, but some mail servers will not send email to your server without it)\n")
	txt.WriteString("2. Create a non-admin matrix user with username `email2matrix` and secure password\n")
	txt.WriteString("3. Create un-encrypted matrix rooms for mailboxes you want to have (1 room = 1 mailbox) and invite email2matrix user in all of them\n")
	txt.WriteString("4. Login with email2matrix user and accept the invites you sent (yes, that's why I asked you if you really want it)\n")
	txt.WriteString("5. Send email2matrix password and pairs of mailbox name - matrix room id to the @support:etke.cc (eg: info@matrix." + o.get("domain") + " = !gqlCuoCdhufltluRXk:" + o.get("domain") + ")\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterEtherpad() string {
	if !o.has("etherpad") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### etherpad\n\n")
	txt.WriteString("1. Open integration manager in any element client app\n")
	txt.WriteString("2. Click on the sprocket icon and go to the Widgets tab\n")
	txt.WriteString("3. Open the Etherpad widget configuration and replace `scalar.vector.im` with `dimension." + o.get("domain") + "`\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("### honoroit\n\n")
	txt.WriteString("1. Create a matrix room (encryption supported) and invite the honoroit user into it\n")
	txt.WriteString("2. Send the room id to the @support:etke.cc\n\n")

	txt.WriteString("if you want to change honoroit messages, send the texts you want to set to the @support:etke.cc (you can use https://gitlab.com/etke.cc/honoroit/-/blob/main/config/defaults.go as reference)\n\n")

	return txt.String()
}
