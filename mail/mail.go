package mail

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/mattevans/postmark-go"
	"github.com/rs/zerolog"
)

// Client to send mail
type Client struct {
	from    string
	replyto string
	sender  *postmark.Client
}

func New(token, from, replyto string) *Client {
	if token == "" {
		return nil
	}
	pm := postmark.NewClient(
		postmark.WithClient(&http.Client{
			Transport: &postmark.AuthTransport{
				Token: token,
			},
		}),
	)

	return &Client{
		from:    from,
		replyto: replyto,
		sender:  pm,
	}
}

func (c *Client) Send(ctx context.Context, req *postmark.Email) error {
	log := zerolog.Ctx(ctx)
	span := sentry.StartSpan(ctx, "http.client", sentry.WithDescription("mail.Send"))
	defer span.Finish()

	req.From = c.from
	req.ReplyTo = c.replyto

	_, resp, err := c.sender.Email.Send(req)
	if err != nil {
		log.Error().Err(err).Any("response", resp).Msg("cannot send email")
		return err
	}

	log.Debug().Str("subject", req.Subject).Str("to", req.To).Msg("email has been sent")
	return nil
}
