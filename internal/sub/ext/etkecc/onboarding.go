package etkecc

import (
	"context"
	"sort"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
	"maunium.net/go/mautrix"
)

func (o *order) generateOnboarding(ctx context.Context) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateOnboarding")
	defer span.Finish()

	log := o.logger(span.Context())
	log.Info().Msg("generating onboarding")

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
			ContentBytes:  []byte(text),
			FileName:      "onboarding.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(text)),
		},
	)
	log.Info().Msg("onboarding has been generated")
}

func (o *order) generateOnboardingIntro() string {
	var txt strings.Builder
	txt.WriteString("Hello!\n\n")
	txt.WriteString("We're happy to announce that your Matrix server on " + o.domain + " is now fully operational and ready for you to use! 🎉\n")
	txt.WriteString("Here's all the essential information you need for a smooth onboarding experience:\n\n")

	return txt.String()
}

func (o *order) generateOnboardingLinks() string {
	var txt strings.Builder

	txt.WriteString("**Server Links**\n\n")
	if !o.has("element-web") {
		txt.WriteString("* Web app: " + link("app.etke.cc") + "\n")
	}
	txt.WriteString("* Homeserver: " + link("matrix."+o.domain) + "\n")
	txt.WriteString("* Synapse Admin: " + link("matrix."+o.domain+"/synapse-admin") + " " + helpLink("etke.cc/help/extras/synapse-admin") + "\n")
	if o.has("etherpad") {
		txt.WriteString("* Etherpad (admin): " + link("etherpad."+o.domain+"/admin") + "\n")
	}
	if o.has("vaultwarden") {
		txt.WriteString("* Vaultwarden (admin):" + link("vault."+o.domain+"/admin") + " " + helpLink("etke.cc/help/extras/vaultwarden") + "\n")
	}
	if o.has("maubot") {
		txt.WriteString("* Maubot (admin): " + link("matrix."+o.domain+"/_matrix/maubot") + "\n")
	}

	items := []string{}
	for item := range dnsmap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, item := range items {
		txt.WriteString("* " + o.c.String(item) + ": " + link(dnsmap[item]+"."+o.domain))
		if helpURL := helpmap[item]; helpURL != "" {
			txt.WriteString(" " + helpLink(helpURL))
		}
		txt.WriteString("\n")
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
	txt.WriteString("**Matrix Bots**\n\n")
	items := []string{}
	for item := range botmap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bot := range items {
		if o.has(bot) {
			txt.WriteString("* " + o.c.String(bot) + ": " + matrixLink(botmap[bot]+":"+o.domain))
			if helpURL := helpmap[bot]; helpURL != "" {
				txt.WriteString(" " + helpLink(helpURL))
			}
			txt.WriteString("\n")
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
	txt.WriteString("**Matrix Bridges**\n\n")
	items := []string{}
	for item := range bridgemap {
		if o.has(item) {
			items = append(items, item)
		}
	}
	sort.Strings(items)
	for _, bridge := range items {
		txt.WriteString("* " + o.c.String(bridge) + ": " + matrixLink(bridgemap[bridge]+":"+o.domain))
		if helpURL := helpmap[bridge]; helpURL != "" {
			txt.WriteString(" " + helpLink(helpURL))
		}
		txt.WriteString("\n")
	}
	txt.WriteString("\n\n")

	return txt.String()
}

func (o *order) generateOnboardingCredentials() string {
	var txt strings.Builder

	// hacky way to print friendly names
	passwords := make(map[string]string)
	for service, password := range o.pass {
		passwords[service] = password
	}
	mxpass := passwords["matrix"]
	delete(o.logins, "matrix")
	delete(passwords, "matrix")
	delete(passwords, "borg")
	delete(passwords, "smtp")

	txt.WriteString("**Matrix Credentials**\n\n")
	txt.WriteString("* Matrix ID: " + matrixLink("@"+o.get("username")+":"+o.domain) + "\n")
	txt.WriteString("* Username: " + o.get("username") + "\n")
	txt.WriteString("* Password: " + mxpass + "\n\n")

	serviceCreds := make([]string, 0, len(o.logins))
	for service := range o.logins {
		serviceCreds = append(serviceCreds, service)
	}
	sort.Strings(serviceCreds)

	for _, service := range serviceCreds {
		txt.WriteString("**" + o.c.String(service) + " Credentials**\n\n")
		txt.WriteString("* Username: " + o.login(service) + "\n")
		txt.WriteString("* Password: " + passwords[service] + "\n\n")
		delete(passwords, service)
	}

	passwordsLeft := []string{}
	for item := range passwords {
		passwordsLeft = append(passwordsLeft, item)
	}
	sort.Strings(passwordsLeft)
	for _, name := range passwordsLeft {
		txt.WriteString("**" + o.c.String(name) + " Credentials**\n\n")
		txt.WriteString("* Username: " + o.get("username") + "\n")
		txt.WriteString("* Password: " + passwords[name] + "\n\n")
	}

	txt.WriteString("Should you encounter any issues or require assistance, please don't hesitate to check out " + link("etke.cc/help") + ".")
	txt.WriteString("We're committed to providing you with the support you need.\n\n")
	return txt.String()
}

func (o *order) generateOnboardingAfter() string {
	var txt strings.Builder

	txt.WriteString(o.generateOnboardingAfterBorgBackup())
	txt.WriteString(o.generateOnboardingAfterBuscarron())
	txt.WriteString(o.generateOnboardingAfterHonoroit())
	txt.WriteString(o.generateOnboardingAfterMigadu())

	text := txt.String()
	if text == "" {
		return ""
	}

	return "**Post-Setup Steps for Specific Components:**\n\n" + text
}

func (o *order) generateOnboardingAfterBorgBackup() string {
	if !o.has("borg") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**borg backup**\n\n")
	txt.WriteString("Add the following public ssh key to the repository configuration on the borg provider side:\n\n")
	txt.WriteString("```\n")
	txt.WriteString(o.password("borg") + "\n")
	txt.WriteString("```\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterBuscarron() string {
	if !o.has("buscarron") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**buscarron**\n\n")
	txt.WriteString("1. Create (encrypted) Matrix room(-s) for the forms you want to have and invite the buscarron user to all of them\n")
	txt.WriteString("2. Send the list of (form name, Matrix room id, redirect URL after submission) to " + matrixLink("@support:etke.cc") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**honoroit**\n\n")
	txt.WriteString("1. Create a matrix room (encryption supported) and invite the honoroit user into it\n")
	txt.WriteString("2. Send the room id to " + matrixLink("@support:etke.cc") + "\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterMigadu() string {
	if !o.has("service-email") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**Email Service**\n\n")
	txt.WriteString("1. Soon, you will receive an invitation to [Migadu](https://migadu.com/)'s admin panel. Check your mailbox\n")
	txt.WriteString("2. Accept the invitation and prepare your account (set a password, etc.)\n")
	txt.WriteString("3. Go to the " + link("admin.migadu.com/domains") + " page and select your email domain\n")
	txt.WriteString("4. Click on **DNS Configuration** and then on **Setup Instructions**\n")
	txt.WriteString("5. Follow the instructions and configure the DNS records for your domain\n")
	txt.WriteString("6. Once that's done, click on **Diagnostics** and check if everything is configured correctly\n\n")

	return txt.String()
}

func (o *order) generateOnboardingOutro() string {
	var txt strings.Builder

	txt.WriteString("You're invited to join the [#news:etke.cc](https://matrix.to/#/#news:etke.cc) room, where we regularly post updates and modifications **pertaining to your server** every week. Please, remain in this room and stay vigilant - each weekly message will include crucial information that we highly recommend you stay informed about. Alternatively, you can also access updates via [Web](https://etke.cc/news/), [RSS](https://etke.cc/news/index.xml), or [Fediverse](https://mastodon.matrix.org/@etkecc).\n\n")

	txt.WriteString("Happy Matrixing!\n\n")

	txt.WriteString("PS: To enrich your Matrix experience and discover vibrant communities, we recommend using ")
	txt.WriteString("[MatrixRooms.info](https://matrixrooms.info?utm_source=etke.cc&utm_medium=email&utm_campaign=onboarding), ")
	txt.WriteString("our own Matrix rooms search engine. It's a fantastic resource for finding and joining interesting Matrix communities.\n\n")

	txt.WriteString("Best regards,\n\n")
	txt.WriteString("etke.cc")

	return txt.String()
}

func matrixLink(id string) string {
	return "[" + id + "](https://matrix.to/#/" + id + ")"
}

func link(label string, address ...string) string {
	if len(address) == 0 {
		return "[" + label + "](https://" + label + ")"
	}

	return "[" + label + "](https://" + address[0] + ")"
}

func helpLink(address string) string {
	return "(" + link("help", address) + ")"
}
