package etkecc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/buscarron/utils"
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
3. Set the custom price to $%d.
4. Subscribe on Ko-Fi with the same email address you used for this order (%s).

Once your payment is confirmed, we'll promptly initiate the setup of your Matrix server. Look forward to a new email that will guide you through the onboarding process with all the necessary details.`
	followupFooter = `
To check the status of your order and stay updated, please keep an eye on your [Order Status Page](%s).

Got any questions? Feel free to reply to this email - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Best regards,

etke.cc`
)

func (o *order) generateFollowup(ctx context.Context, questions, delegation, dns string, countQ int) {
	span := utils.StartSpan(ctx, "sub.ext.etkecc.generateFollowup")
	defer span.Finish()

	var txt strings.Builder
	txt.WriteString(followupHeader)
	if countQ > 0 {
		txt.WriteString(questions)
	} else {
		txt.WriteString(fmt.Sprintf(followupNoQuestions, o.price, o.get("email")))
	}

	if o.hosting == "" {
		txt.WriteString("\n\n")
		txt.WriteString("Please, ensure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls).")

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

	content := format.RenderMarkdown(txt.String(), true, true)
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
}

func (o *order) sendFollowup(ctx context.Context) {
	if o.pm == nil || (reflect.ValueOf(o.pm).Kind() == reflect.Ptr && reflect.ValueOf(o.pm).IsNil()) {
		return
	}

	req := &postmark.Email{
		To:            o.get("email"),
		Tag:           "confirmation",
		Subject:       "Your Matrix Server Order Update (" + o.domain + ") 🚀",
		MessageStream: "followups",
		TextBody:      o.followup.Body,
		HTMLBody:      o.followup.FormattedBody,
	}
	o.pm.Send(ctx, req) //nolint:errcheck // no need to wait
}
