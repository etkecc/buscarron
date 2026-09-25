package mail

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type smtpSuite struct {
	suite.Suite
}

func (s *smtpSuite) TestBuildSMTPMessage() {
	tests := []struct {
		name    string
		email   *postmark.Email
		errored bool
		check   func(s *suite.Suite, msg string)
	}{
		{
			name: "text only",
			email: &postmark.Email{
				From:     "from@example.com",
				To:       "to@example.com",
				Subject:  "hello",
				TextBody: "text body",
			},
			check: func(s *suite.Suite, msg string) {
				s.Contains(msg, "From: from@example.com\r\n")
				s.Contains(msg, "To: to@example.com\r\n")
				s.Contains(msg, "Subject: hello\r\n")
				s.Contains(msg, "MIME-Version: 1.0\r\n")
				s.Contains(msg, "Content-Type: text/plain; charset=utf-8\r\n")
				s.Contains(msg, "\r\ntext body")
			},
		},
		{
			name: "html only",
			email: &postmark.Email{
				From:     "from@example.com",
				To:       "to@example.com",
				Subject:  "hello",
				HTMLBody: "<p>html body</p>",
			},
			check: func(s *suite.Suite, msg string) {
				s.Contains(msg, "Content-Type: text/html; charset=utf-8\r\n")
				s.Contains(msg, "\r\n<p>html body</p>")
			},
		},
		{
			name: "multipart",
			email: &postmark.Email{
				From:     "from@example.com",
				To:       "to@example.com",
				Subject:  "hello",
				TextBody: "text part",
				HTMLBody: "<p>html part</p>",
			},
			check: func(s *suite.Suite, msg string) {
				s.Contains(msg, "Content-Type: multipart/alternative; boundary=buscarron-alternative\r\n")
				s.Contains(msg, "--buscarron-alternative\r\nContent-Type: text/plain; charset=utf-8\r\n\r\ntext part")
				s.Contains(msg, "--buscarron-alternative\r\nContent-Type: text/html; charset=utf-8\r\n\r\n<p>html part</p>")
				s.Contains(msg, "--buscarron-alternative--\r\n")
			},
		},
		{
			name: "optional headers",
			email: &postmark.Email{
				From:     "from@example.com",
				To:       "to@example.com",
				Cc:       "cc@example.com",
				Bcc:      "bcc@example.com",
				ReplyTo:  "reply@example.com",
				Tag:      "confirmation",
				Subject:  "hello",
				TextBody: "text body",
			},
			check: func(s *suite.Suite, msg string) {
				s.Contains(msg, "Cc: cc@example.com\r\n")
				s.NotContains(msg, "Bcc:")
				s.Contains(msg, "Reply-To: reply@example.com\r\n")
				s.Contains(msg, "X-Tag: confirmation\r\n")
			},
		},
		{
			name: "crlf in header rejected",
			email: &postmark.Email{
				From:    "from@example.com",
				To:      "to@example.com\r\nBcc: evil@example.com",
				Subject: "hello",
			},
			errored: true,
		},
		{
			name: "non-ascii subject",
			email: &postmark.Email{
				From:     "from@example.com",
				To:       "to@example.com",
				Subject:  "Update 🚀",
				TextBody: "text body",
			},
			check: func(s *suite.Suite, msg string) {
				subject := subjectLine(msg)
				s.Contains(subject, "?b?")
				s.NotContains(msg, "Update 🚀")
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			msg, err := buildSMTPMessage(tt.email)
			if tt.errored {
				s.Error(err)
				return
			}
			s.NoError(err)
			tt.check(&s.Suite, string(msg))
		})
	}
}

func (s *smtpSuite) TestSMTPSend() {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	s.Require().NoError(err)
	defer l.Close()

	var data bytes.Buffer
	done := make(chan error, 1)
	go func() {
		done <- s.serveSMTP(l, &data)
	}()

	na := l.Addr()
	addr, ok := na.(*net.TCPAddr)
	s.Require().True(ok)
	port := strconv.Itoa(addr.Port)
	smtpc := NewSMTP("127.0.0.1", port, "", "")
	err = smtpc.Send(s.T().Context(), &postmark.Email{
		From:     "from@example.com",
		To:       "to@example.com",
		TextBody: "hello",
	})
	s.Require().NoError(err)
	s.Require().NoError(<-done)

	msg := data.String()
	s.Contains(msg, "From: from@example.com")
	s.Contains(msg, "To: to@example.com")
	s.Contains(msg, "hello")
}

func (s *smtpSuite) serveSMTP(l net.Listener, data *bytes.Buffer) error {
	conn, err := l.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	_, _ = conn.Write([]byte("220 fake ESMTP\r\n"))
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimRight(line, "\r\n")
		switch {
		case strings.HasPrefix(line, "EHLO"):
			_, _ = conn.Write([]byte("250-fake\r\n250 SIZE 1000000\r\n"))
		case strings.HasPrefix(line, "MAIL FROM"):
			_, _ = conn.Write([]byte("250 ok\r\n"))
		case strings.HasPrefix(line, "RCPT TO"):
			_, _ = conn.Write([]byte("250 ok\r\n"))
		case strings.HasPrefix(line, "DATA"):
			_, _ = conn.Write([]byte("354 go ahead\r\n"))
			if err := readUntilDot(r, data); err != nil {
				return err
			}
			_, _ = conn.Write([]byte("250 ok\r\n"))
		case strings.HasPrefix(line, "QUIT"):
			_, _ = conn.Write([]byte("221 bye\r\n"))
			return nil
		default:
			_, _ = conn.Write([]byte("250 ok\r\n"))
		}
	}
}

func readUntilDot(r *bufio.Reader, data *bytes.Buffer) error {
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if line == ".\r\n" {
			return nil
		}
		data.WriteString(line)
	}
}

func subjectLine(msg string) string {
	for _, line := range strings.Split(msg, "\r\n") {
		if strings.HasPrefix(line, "Subject: ") {
			return line
		}
	}
	return ""
}

func TestSMTP(t *testing.T) {
	suite.Run(t, new(smtpSuite))
}
