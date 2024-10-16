package etkecc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/etkecc/buscarron/internal/utils"
	formatCustom "github.com/etkecc/go-kit/format"
	"github.com/mattevans/postmark-go"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/format"
)

const (
	followupHeader = `Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

`
	followupNoQuestions = `No need for additional details at this moment; we'll keep it simple:

**Payment Instructions**:

1. Visit our [membership page](https://etke.cc/membership).
2. Select the "By Complexity" tier.
3. Set the custom price to **$%d**.
4. Subscribe on Ko-Fi with the same email address you used for this order (%s).

Once your payment is confirmed, we'll promptly initiate the setup of your Matrix server. Look forward to a new email that will guide you through the onboarding process with all the necessary details.`
	followupFooter = `
To check the status of your order and stay updated, please keep an eye on your [Order Status Page](%s).

Got any questions? Feel free to [contact us](https://etke.cc/contacts/) - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team`
)

func (o *order) generateFollowup(ctx context.Context, questions, delegation, dns string, countQ int) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateFollowup")
	defer span.Finish()

	log := o.logger(ctx)
	log.Info().Msg("generating followup")

	var txt strings.Builder
	txt.WriteString(followupHeader)
	if countQ > 0 {
		txt.WriteString(questions)
	} else {
		txt.WriteString(fmt.Sprintf(followupNoQuestions, o.price, o.get("email")))
	}

	if o.hosting == "" {
		authorizedkeys := "/home/" + o.get("ssh-user") + "/.ssh/authorized_keys"
		if o.get("ssh-user") == "root" {
			authorizedkeys = "/root/.ssh/authorized_keys"
		}

		pubkey := "technical issue; the key will be provided later, we apologize for the inconvenience"
		for _, file := range o.files {
			if file.FileName == "sshkey.pub" {
				pubkey = string(file.ContentBytes)
				break
			}
		}

		txt.WriteString("\n\n")
		txt.WriteString("Please, ensure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), ")
		txt.WriteString("and the following ssh key added to the **" + authorizedkeys + "**: `" + pubkey + "`. ")
		txt.WriteString("Here is how you can do that:\n\n")
		txt.WriteString("1. ensure the .ssh directory exists: `mkdir -p " + strings.TrimSuffix(authorizedkeys, "/authorized_keys") + "`\n")
		txt.WriteString("2. ensure the authorized_keys file exists: `touch " + authorizedkeys + "`\n")
		txt.WriteString("3. add the key to the authorized_keys file: `echo '" + pubkey + "' >> " + authorizedkeys + "`\n")
		txt.WriteString("4. ensure the correct permissions are set on the authorized_keys file: `chmod 600 " + authorizedkeys + "`\n")
		txt.WriteString("5. ensure the correct permissions are set on the .ssh directory: `chmod 700 " + strings.TrimSuffix(authorizedkeys, "/authorized_keys") + "`\n")
		txt.WriteString("6. ensure the correct ownership is set on the .ssh directory: `chown -hR " + o.get("ssh-user") + ":" + o.get("ssh-user") + " " + strings.TrimSuffix(authorizedkeys, "/authorized_keys") + "`\n")

		if dns != "" {
			txt.WriteString("\n\n")
			txt.WriteString(dns)
		}
	}

	if delegation != "" {
		txt.WriteString("\n\n")
		txt.WriteString(delegation)
	}

	h := sha256.New()
	h.Write([]byte(o.domain))
	id := hex.EncodeToString(h.Sum(nil))
	txt.WriteString(fmt.Sprintf(followupFooter, "https://etke.cc/order/status/#"+id))

	text := txt.String()
	o.response = formatCustom.Render(text)
	content := format.RenderMarkdown(text, true, true)
	o.followup = &content
	o.files = append(o.files,
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(content.Body),
			ContentBytes:  []byte(content.Body),
			FileName:      "followup.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(content.Body)),
		},
	)
	log.Info().Msg("followup has been generated")
}

func (o *order) sendFollowup(ctx context.Context) {
	ctx = context.WithoutCancel(ctx)
	if o.pm == nil || (reflect.ValueOf(o.pm).Kind() == reflect.Ptr && reflect.ValueOf(o.pm).IsNil()) {
		return
	}

	log := o.logger(ctx)
	log.Info().Msg("sending followup")

	req := &postmark.Email{
		To:            o.get("email"),
		Tag:           "confirmation",
		Subject:       "Your Matrix Server Order Update (" + o.domain + ") 🚀",
		MessageStream: "followups",
		TextBody:      o.followup.Body,
		HTMLBody:      o.followup.FormattedBody,
	}
	if err := o.pm.Send(ctx, req); err != nil {
		log.Error().Err(err).Msg("cannot send followup")
	}
	log.Info().Msg("followup has been sent")
}
