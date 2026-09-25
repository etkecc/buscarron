package mail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mattevans/postmark-go"
)

// backendURL of the postmark api, test hook
var backendURL = "https://api.postmarkapp.com"

// PostmarkSender sends mail via postmark
type PostmarkSender struct {
	client *postmark.Client
}

// NewPostmark sender
func NewPostmark(token string) Sender {
	return &PostmarkSender{client: postmark.NewClient(
		postmark.WithClient(&http.Client{Transport: &postmark.AuthTransport{Token: token}}),
		postmark.WithBackendURL(backendURL),
	)}
}

// Send email via postmark
func (p *PostmarkSender) Send(_ context.Context, req *postmark.Email) error {
	_, resp, err := p.client.Email.Send(req)
	if err != nil {
		return fmt.Errorf("postmark: %w (status: %d, message: %s)", err, resp.StatusCode, resp.Message)
	}

	return nil
}
