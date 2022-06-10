package mail

import (
	"net/http"

	"github.com/mattevans/postmark-go"
	"gitlab.com/etke.cc/go/logger"
)

// Client to send mail
type Client struct {
	log     *logger.Logger
	from    string
	replyto string
	sender  *postmark.Client
}

func New(token, from, replyto, loglevel string) *Client {
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
		log:     logger.New("mail.", loglevel),
	}
}

func (c *Client) Send(req *postmark.Email) error {
	req.From = c.from
	req.ReplyTo = c.replyto

	_, resp, err := c.sender.Email.Send(req)
	if err != nil {
		c.log.Error("cannot send email: %v (response: %+v)", err, resp)
		return err
	}

	c.log.Debug("email '%s' was sent to %s", req.Subject, req.To)
	return nil
}
