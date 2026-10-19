package mail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/etkecc/go-kit/httpclient"
	"github.com/mattevans/postmark-go"
)

// PostmarkSender sends mail via the Postmark API
type PostmarkSender struct {
	client *postmark.Client
}

// NewPostmark creates the Postmark transport, opts are for tests
func NewPostmark(token string, opts ...postmark.Option) Sender {
	// never follow redirects, the token must not leak to another host
	httpClient := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport:     &postmark.AuthTransport{Token: token, Transport: httpclient.NewTransport()},
	}
	options := append([]postmark.Option{postmark.WithClient(httpClient)}, opts...)
	client := postmark.NewClient(options...)

	return &PostmarkSender{client: client}
}

// Send sends the email via the Postmark API
func (p *PostmarkSender) Send(_ context.Context, req *postmark.Email) error {
	_, resp, err := p.client.Email.Send(req)
	if err == nil {
		return nil
	}
	if resp == nil {
		return fmt.Errorf("postmark: %w", err)
	}

	return fmt.Errorf("postmark: %w (status: %d)", err, resp.StatusCode)
}
