package mail //nolint:var-naming // Package mail provides a client to send mail via a configured sender

import (
	"context"

	"github.com/mattevans/postmark-go"
	"github.com/rs/zerolog"
)

// Sender is a mail transport
type Sender interface {
	Send(ctx context.Context, req *postmark.Email) error
}

// Client stamps the email with sender fields and delegates to the transport
type Client struct {
	from    string
	replyto string
	sender  Sender
}

// New creates the mail client, nil sender yields nil client
func New(from, replyto string, sender Sender) *Client {
	if sender == nil {
		return nil
	}

	return &Client{
		from:    from,
		replyto: replyto,
		sender:  sender,
	}
}

// Send sends the email
func (c *Client) Send(ctx context.Context, req *postmark.Email) error {
	log := zerolog.Ctx(ctx)

	req.From = c.from
	req.ReplyTo = c.replyto

	if err := c.sender.Send(ctx, req); err != nil {
		log.Error().Err(err).Msg("cannot send email")
		return err
	}

	log.Debug().Str("subject", req.Subject).Str("to", req.To).Msg("email has been sent")
	return nil
}
