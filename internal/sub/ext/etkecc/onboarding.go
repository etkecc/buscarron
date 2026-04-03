package etkecc

import (
	"context"
	"fmt"
	"maps"
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
	txt.WriteString("Great news — your Matrix server on " + o.domain + " is up and running and ready to use! 🎉\n")
	txt.WriteString("Below is everything you need to get started:\n\n")

	return txt.String()
}

func (o *order) generateOnboardingLinks() string {
	var txt strings.Builder

	txt.WriteString("**Server Links**\n\n")
	if !o.has("element-web") {
		txt.WriteString("* Web app: " + link("app.etke.cc") + "\n")
	}
	txt.WriteString("* Homeserver: " + link("matrix."+o.domain) + "\n")
	ketesaLink := fmt.Sprintf("[matrix.%s/admin](https://matrix.%s/admin/?username=%s)", o.domain, o.domain, o.get("username"))
	txt.WriteString("* Admin Panel: " + ketesaLink + " " + helpLink("etke.cc/help/extras/ketesa") + "\n")
	if o.has("service-email") {
		txt.WriteString("* Email Service (admin): " + link("admin.migadu.com") + "\n")
	}
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
	maps.Copy(passwords, o.pass)
	mxpass := passwords["matrix"]
	delete(o.logins, "matrix")
	delete(passwords, "matrix")
	delete(passwords, "borg")
	delete(passwords, "smtp")

	txt.WriteString("**Matrix Credentials**\n\n")
	txt.WriteString("* Matrix ID: " + matrixLink("@"+o.get("username")+":"+o.domain) + "\n")
	txt.WriteString("* Username: " + o.get("username") + "\n")
	txt.WriteString("* Password: " + mxpass + "\n\n")

	if o.has("radicale") {
		txt.WriteString("**Radicale (CalDAV/CardDAV) Credentials**\n\n")
		txt.WriteString("* Username: your Matrix username (" + o.get("username") + ")\n")
		txt.WriteString("* Password: your Matrix password\n\n")
		txt.WriteString("> Radicale can be used by any user of your Matrix server with their Matrix credentials, thanks to [radicale-auth-matrix](https://github.com/etkecc/radicale-auth-matrix).\n\n")
	}

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

	passwordsLeft := make([]string, 0, len(passwords))
	for item := range passwords {
		passwordsLeft = append(passwordsLeft, item)
	}
	sort.Strings(passwordsLeft)
	for _, name := range passwordsLeft {
		txt.WriteString("**" + o.c.String(name) + " Credentials**\n\n")
		txt.WriteString("* Username: " + o.get("username") + "\n")
		txt.WriteString("* Password: " + passwords[name] + "\n\n")
	}

	txt.WriteString("If you run into any issues or have questions, please visit " + link("etke.cc/help") + ". ")
	txt.WriteString("We're here to help.\n\n")
	return txt.String()
}

func (o *order) generateOnboardingAfter() string {
	var txt strings.Builder

	txt.WriteString(o.generateOnboardingAfterBorgBackup())
	txt.WriteString(o.generateOnboardingAfterBuscarron())
	txt.WriteString(o.generateOnboardingAfterDraupnir())
	txt.WriteString(o.generateOnboardingAfterHonoroit())
	txt.WriteString(o.generateOnboardingAfterMigadu())

	text := txt.String()
	if text == "" {
		return ""
	}

	return "**Post-Setup Steps for Specific Components**\n\n" + text
}

func (o *order) generateOnboardingAfterBorgBackup() string {
	if !o.has("borg") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**borg backup**\n\n")
	txt.WriteString("Please add the following public SSH key to the repository configuration in your Borg provider:\n\n")
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
	txt.WriteString("1. Create encrypted Matrix rooms for the forms you want, and invite the buscarron user to each room\n")
	txt.WriteString("2. Send us the list of form name, Matrix room ID, and redirect URL after submission [here](https://etke.cc/contacts/)\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterDraupnir() string {
	if !o.has("draupnir") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**Draupnir** (moderation bot)\n\n")
	txt.WriteString("1. Create an unencrypted private Matrix room to use as the management room\n")
	txt.WriteString("2. Send the room ID [to us](https://etke.cc/contacts/)\n")
	txt.WriteString("3. Once you receive confirmation that the bot is ready, follow [the step-by-step guide](https://etke.cc/help/bots/draupnir/#getting-started) to finish the setup\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterHonoroit() string {
	if !o.has("honoroit") {
		return ""
	}
	var txt strings.Builder

	txt.WriteString("**honoroit**\n\n")
	txt.WriteString("1. Create a Matrix room (encryption supported) and invite the honoroit user\n")
	txt.WriteString("2. Send the room ID [to us](https://etke.cc/contacts/)\n\n")

	return txt.String()
}

func (o *order) generateOnboardingAfterMigadu() string {
	if !o.has("service-email") || o.subdomain {
		return ""
	}
	var txt strings.Builder

	var orderEmailHost string
	orderEmailParts := strings.Split(o.get("email"), "@")
	if len(orderEmailParts) == 2 {
		orderEmailHost = orderEmailParts[1]
	}

	txt.WriteString("**Email Service**\n\n")
	if orderEmailHost == o.domain {
		txt.WriteString("\n\n> **WARNING!**: You have ordered [etke.cc email service](https://etke.cc/help/extras/email-hosting/) for your Matrix server's domain, **BUT you already have an email service** configured on that domain! Configuring the new email service **will break your existing email service!**. Please, ensure you are ready for that!\n\n")
	}
	txt.WriteString("1. You'll receive an invitation to [Migadu](https://migadu.com/)'s admin panel soon. Please check your inbox\n")
	txt.WriteString("2. Accept the invitation and prepare your account (set a password, etc.)\n")
	txt.WriteString("3. Go to the " + link("admin.migadu.com/domains") + " page and select your email domain\n")
	txt.WriteString("4. Click **DNS Configuration**, then **Setup Instructions**\n")
	txt.WriteString("5. Follow the instructions to configure the DNS records for your domain\n")
	txt.WriteString("6. When you're done, click **Diagnostics** and confirm everything is configured correctly\n\n")

	return txt.String()
}

func (o *order) generateOnboardingOutro() string {
	var txt strings.Builder

	txt.WriteString("You're invited to join the [#news:etke.cc](https://matrix.to/#/#news:etke.cc) room, where we post weekly updates and changes **about your server**. Please stay in the room so you don't miss important announcements. You can also read updates via [Web](https://etke.cc/news/), [RSS](https://etke.cc/news/index.xml), or [Fediverse](https://mastodon.matrix.org/@etkecc).\n\n")

	txt.WriteString("Happy Matrixing!\n\n")

	txt.WriteString("PS: Looking for active Matrix communities? Try ")
	txt.WriteString("[MatrixRooms.info](https://matrixrooms.info?utm_source=etke.cc&utm_medium=email&utm_campaign=onboarding), ")
	txt.WriteString("our Matrix room search engine. It's an easy way to discover and join new rooms.\n\n")

	txt.WriteString("Warm regards,\n\n")
	txt.WriteString("the [etke.cc](https://etke.cc) team\n")

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
