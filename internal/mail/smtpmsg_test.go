package mail

import (
	"bufio"
	"bytes"
	"net/textproto"
	"regexp"
	"strings"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type smtpmsgSuite struct {
	suite.Suite
}

func (s *smtpmsgSuite) TestBuildMessage_Text() {
	msg, err := buildMessage(s.req())

	s.NoError(err)
	text := string(msg)
	s.Contains(text, "From: from@example.com\r\n")
	s.Contains(text, "To: to@example.com\r\n")
	s.Contains(text, "Subject: subject\r\n")
	s.Contains(text, "Date: ")
	s.Contains(text, "MIME-Version: 1.0\r\n")
	s.Contains(text, "Content-Type: text/plain; charset=utf-8\r\n\r\ntext body")
	s.Contains(text, "Message-ID: <")
	s.Contains(text, "@example.com>")
}

func (s *smtpmsgSuite) TestBuildMessage_HTML() {
	req := s.req()
	req.TextBody = ""
	req.HTMLBody = "<p>html body</p>"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Contains(string(msg), "Content-Type: text/html; charset=utf-8\r\n\r\n<p>html body</p>")
}

func (s *smtpmsgSuite) TestBuildMessage_Alternative() {
	req := s.req()
	req.HTMLBody = "<p>html body</p>"
	msg, err := buildMessage(req)

	s.NoError(err)
	text := string(msg)
	boundary := s.boundary(text)
	s.Contains(text, `Content-Type: multipart/alternative; boundary="`+boundary+`"`)
	s.Contains(text, "--"+boundary+"\r\nContent-Type: text/plain; charset=utf-8\r\n\r\ntext body\r\n")
	s.Contains(text, "--"+boundary+"\r\nContent-Type: text/html; charset=utf-8\r\n\r\n<p>html body</p>\r\n")
	s.Contains(text, "--"+boundary+"--\r\n")
	s.NotContains(text, "multipart/mixed")
}

func (s *smtpmsgSuite) TestBuildMessage_Tag() {
	req := s.req()
	req.Tag = "confirmation"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Contains(string(msg), "X-Tag: confirmation\r\n")
	s.NotContains(string(msg), "X-Tag: \r\n")
}

func (s *smtpmsgSuite) TestBuildMessage_Subject() {
	req := s.req()
	req.Subject = "Update 🚀"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Contains(string(msg), "Subject: =?utf-8?b?")
	s.NotContains(string(msg), "🚀")
}

func (s *smtpmsgSuite) TestBuildMessage_Bcc() {
	req := s.req()
	req.Bcc = "bcc@example.com"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.NotContains(string(msg), "bcc@example.com")
	s.NotContains(string(msg), "Bcc")
}

func (s *smtpmsgSuite) TestBuildMessage_LF() {
	req := s.req()
	req.TextBody = "line1\nline2"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Contains(string(msg), "line1\r\nline2")
	s.NotContains(string(msg), "line1\nline2")
}

func (s *smtpmsgSuite) TestBuildMessage_LoneCR() {
	req := s.req()
	req.TextBody = "A\r.\rB"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.NotContains(string(msg), "\r.")
	var wire bytes.Buffer
	writer := textproto.NewWriter(bufio.NewWriter(&wire)).DotWriter()
	_, err = writer.Write(msg)
	s.NoError(err)
	s.NoError(writer.Close())
	s.Contains(wire.String(), "\r\n..\r\n")
}

func (s *smtpmsgSuite) TestBuildMessage_ASCIIBody() {
	msg, err := buildMessage(s.req())

	s.NoError(err)
	s.NotContains(string(msg), "Content-Transfer-Encoding")
}

func (s *smtpmsgSuite) TestBuildMessage_NonASCIIBody() {
	req := s.req()
	req.TextBody = "Hello 🎉"
	req.HTMLBody = "<p>Hello 🎉</p>"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Equal(2, strings.Count(string(msg), "Content-Transfer-Encoding: 8bit\r\n"))
}

func (s *smtpmsgSuite) TestBuildMessage_MixedBody() {
	req := s.req()
	req.HTMLBody = "<p>Hello 🎉</p>"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.Equal(1, strings.Count(string(msg), "Content-Transfer-Encoding: 8bit\r\n"))
}

func (s *smtpmsgSuite) TestBuildMessage_BccOnly() {
	req := s.req()
	req.To = ""
	req.Bcc = "bcc@example.com"
	msg, err := buildMessage(req)

	s.NoError(err)
	s.NotContains(string(msg), "To: ")
	s.NotContains(string(msg), "bcc@example.com")
}

func (s *smtpmsgSuite) TestBuildMessage_Errors() {
	s.Run("empty from", func() {
		req := s.req()
		req.From = ""
		_, err := buildMessage(req)
		s.EqualError(err, "mail: empty from address")
	})
	s.Run("empty body", func() {
		req := s.req()
		req.TextBody = ""
		_, err := buildMessage(req)
		s.EqualError(err, "mail: empty message body")
	})
}

func (s *smtpmsgSuite) TestBuildMessage_Injection() {
	tests := []struct {
		name  string
		apply func(*postmark.Email)
		err   string
	}{
		{"from", func(r *postmark.Email) { r.From = "from@example.com\r\nBcc: evil@example.com" }, "line breaks"},
		{"to", func(r *postmark.Email) { r.To = "to@example.com\r\nBcc: evil@example.com" }, "line breaks"},
		{"cc", func(r *postmark.Email) { r.Cc = "cc@example.com\r\nBcc: evil@example.com" }, "line breaks"},
		{"bcc", func(r *postmark.Email) { r.Bcc = "bcc@example.com\r\nBcc: evil@example.com" }, "line breaks"},
		{"replyto", func(r *postmark.Email) { r.ReplyTo = "reply@example.com\r\nBcc: evil@example.com" }, "line breaks"},
		{"tag", func(r *postmark.Email) { r.Tag = "tag\r\nBcc: evil@example.com" }, "line breaks"},
		{"subject", func(r *postmark.Email) { r.Subject = "subject\r\nBcc: evil@example.com" }, "line breaks"},
		{"header value", func(r *postmark.Email) {
			r.Headers = []postmark.EmailHeader{{Name: "X-Custom", Value: "value\r\nBcc: evil@example.com"}}
		}, "line breaks"},
		{"header name breaks", func(r *postmark.Email) {
			r.Headers = []postmark.EmailHeader{{Name: "X-Custom\r\nBcc: evil@example.com", Value: "value"}}
		}, "line breaks"},
		{"header name colon", func(r *postmark.Email) {
			r.Headers = []postmark.EmailHeader{{Name: "X-Custom: evil", Value: "value"}}
		}, "colon"},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.req()
			tt.apply(req)
			_, err := buildMessage(req)
			s.ErrorContains(err, tt.err)
		})
	}
}

func (s *smtpmsgSuite) TestParseRecipients() {
	recipients, err := parseRecipients(&postmark.Email{
		To: "To Name <to@example.com>", Cc: "cc@example.com", Bcc: "bcc@example.com",
	})

	s.NoError(err)
	s.Equal([]string{"to@example.com", "cc@example.com", "bcc@example.com"}, recipients)
}

func (s *smtpmsgSuite) TestParseRecipients_CcOnly() {
	recipients, err := parseRecipients(&postmark.Email{Cc: "cc@example.com"})

	s.NoError(err)
	s.Equal([]string{"cc@example.com"}, recipients)
}

func (s *smtpmsgSuite) TestParseRecipients_Empty() {
	_, err := parseRecipients(&postmark.Email{})

	s.EqualError(err, "mail: no recipients")
}

func (s *smtpmsgSuite) TestParseRecipients_Group() {
	_, err := parseRecipients(&postmark.Email{To: "undisclosed-recipients:;"})

	s.EqualError(err, "mail: no recipients")
}

func (s *smtpmsgSuite) TestParseRecipients_Invalid() {
	_, err := parseRecipients(&postmark.Email{To: "not an email"})

	s.ErrorContains(err, "mail: cannot parse recipients")
}

func (s *smtpmsgSuite) req() *postmark.Email {
	return &postmark.Email{From: "from@example.com", To: "to@example.com", Subject: "subject", TextBody: "text body"}
}

func (s *smtpmsgSuite) boundary(msg string) string {
	match := regexp.MustCompile(`boundary="([^"]+)"`).FindStringSubmatch(msg)
	s.Require().Len(match, 2)

	return match[1]
}

func TestSMTPMsg(t *testing.T) {
	suite.Run(t, new(smtpmsgSuite))
}
