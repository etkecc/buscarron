package mail //nolint:var-naming // Package mail provides a client to send mail

import (
	"context"

	"github.com/mattevans/postmark-go"
	"github.com/rs/zerolog"
)

// Sender sends mail
type Sender interface {
	Send(ctx context.Context, req *postmark.Email) error
}

// Client to send mail
type Client struct {
	from    string
	replyto string
	sender  Sender
}

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

func (c *Client) Send(ctx context.Context, req *postmark.Email) error {
	log := zerolog.Ctx(ctx)

	req.From = c.from
	req.ReplyTo = c.replyto

	err := c.sender.Send(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("cannot send email")
		return err
	}

	log.Debug().Str("subject", req.Subject).Str("to", req.To).Msg("email has been sent")
	return nil
}
