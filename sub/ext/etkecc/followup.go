package etkecc

import (
	"fmt"
	"reflect"
	"strings"

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
3. Set the custom price to $%d.
4. Use the email you've sent within the order for your Ko-Fi subscription

Once your payment is confirmed, we'll promptly initiate the setup of your Matrix server. Look out for an upcoming email that will guide you through the onboarding process with all the necessary details.`
	followupFooter = `
Got any questions? Feel free to reply to this email; we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Best regards,

Buscarron Stacks,
etke.cc`
)

func (o *order) generateFollowup(questions, dns string, countQ int, dnsInternal bool) {
	var txt strings.Builder
	txt.WriteString(followupHeader)
	if countQ > 0 {
		txt.WriteString(questions)
	} else {
		txt.WriteString(fmt.Sprintf(followupNoQuestions, o.price))
	}

	if o.hosting == "" && !dnsInternal {
		txt.WriteString("\n")
		txt.WriteString(dns)
	}

	txt.WriteString(followupFooter)

	content := format.RenderMarkdown(txt.String(), true, true)
	o.followup = &content
	o.files = append(o.files,
		&mautrix.ReqUploadMedia{
			Content:       strings.NewReader(content.Body),
			FileName:      "onboarding.md",
			ContentType:   "text/markdown",
			ContentLength: int64(len(content.Body)),
		},
	)
}

func (o *order) sendFollowup() {
	if o.pm == nil || (reflect.ValueOf(o.pm).Kind() == reflect.Ptr && reflect.ValueOf(o.pm).IsNil()) {
		return
	}

	req := &postmark.Email{
		To:       o.get("email"),
		Tag:      "confirmation",
		Subject:  "Your Matrix Server Order Update 🚀",
		TextBody: o.followup.Body,
		HTMLBody: o.followup.FormattedBody,
	}
	err := o.pm.Send(req)
	if err != nil {
		o.txt.WriteString("\n\nfollowup: ❌\n")
		return
	}

	o.txt.WriteString("\n\nfollowup: ✅\n")
}
