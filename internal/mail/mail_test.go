package mail

import (
	"context"
	"errors"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type fakeSender struct {
	req *postmark.Email
	err error
}

func (f *fakeSender) Send(_ context.Context, req *postmark.Email) error {
	f.req = req
	return f.err
}

type mailSuite struct {
	suite.Suite
}

func (s *mailSuite) TestNew() {
	client := New("from@example.com", "reply@example.com", NewPostmark("token"))

	s.IsType(&Client{}, client)
}

func (s *mailSuite) TestNew_Empty() {
	s.Nil(New("", "", nil))
}

func (s *mailSuite) TestClient_Send() {
	fake := &fakeSender{}
	client := New("from@example.com", "reply@example.com", fake)
	req := &postmark.Email{To: "to@example.com", Subject: "subject", TextBody: "body"}

	s.NoError(client.Send(s.T().Context(), req))
	s.Same(req, fake.req)
	s.Equal("from@example.com", req.From)
	s.Equal("reply@example.com", req.ReplyTo)

	fake.err = errors.New("boom")
	s.EqualError(client.Send(s.T().Context(), req), "boom")
}

func TestMail(t *testing.T) {
	suite.Run(t, new(mailSuite))
}
