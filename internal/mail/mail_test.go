package mail

import (
	"context"
	"errors"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type mailSuite struct {
	suite.Suite
}

type fakeSender struct {
	req *postmark.Email
	err error
}

func (f *fakeSender) Send(_ context.Context, req *postmark.Email) error {
	f.req = req
	return f.err
}

func (s *mailSuite) TestNew() {
	pm := New("from@example.com", "reply@example.com", NewPostmark("token"))

	s.IsType(&Client{}, pm)
}

func (s *mailSuite) TestNew_Empty() {
	null := New("from@example.com", "reply@example.com", nil)

	s.Nil(null)
}

func (s *mailSuite) TestClient_Send() {
	fs := &fakeSender{}
	client := New("from@example.com", "reply@example.com", fs)
	req := &postmark.Email{To: "to@example.com", Subject: "s", TextBody: "b"}

	s.NoError(client.Send(s.T().Context(), req))
	s.Equal(req, fs.req)
	s.Equal("from@example.com", req.From)
	s.Equal("reply@example.com", req.ReplyTo)

	fs.err = errors.New("boom")
	s.EqualError(client.Send(s.T().Context(), req), "boom")
}

func TestMail(t *testing.T) {
	suite.Run(t, new(mailSuite))
}
