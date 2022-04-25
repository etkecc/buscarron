package mail

import (
	"net/http"

	"github.com/mattevans/postmark-go"
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
	pm := postmark.NewClient(&http.Client{
		Transport: &postmark.AuthTransport{
			Token: token,
		},
	})

	return &Client{
		from:    from,
		replyto: replyto,
		sender:  pm,
	}
}

func (c *Client) Send(req *postmark.Email) (*postmark.EmailResponse, *postmark.Response, error) {
	req.From = c.from
	req.ReplyTo = c.replyto

	return c.sender.Email.Send(req)
}
