package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mattevans/postmark-go"
)

// boundary of the multipart message, fixed to keep message construction deterministic
const boundary = "buscarron-alternative"

// SMTPSender sends mail via smtp
type SMTPSender struct {
	host     string
	port     string
	login    string
	password string
}

// NewSMTP sender
func NewSMTP(host, port, login, password string) Sender {
	return &SMTPSender{host: host, port: port, login: login, password: password}
}

// Send email via smtp
func (s *SMTPSender) Send(ctx context.Context, req *postmark.Email) error {
	msg, err := buildSMTPMessage(req)
	if err != nil {
		return err
	}

	recipients := []string{req.To}
	if req.Cc != "" {
		recipients = append(recipients, req.Cc)
	}
	if req.Bcc != "" {
		recipients = append(recipients, req.Bcc)
	}

	return s.deliver(ctx, req.From, recipients, msg)
}

// authenticate performs plain auth when login is set
func (s *SMTPSender) authenticate(c *smtp.Client) error {
	if s.login == "" {
		return nil
	}
	if ok, _ := c.Extension("AUTH"); !ok {
		return fmt.Errorf("mail: smtp server does not advertise AUTH")
	}
	return c.Auth(smtp.PlainAuth("", s.host, s.login, s.password))
}

func (s *SMTPSender) deliver(ctx context.Context, from string, recipients []string, msg []byte) error {
	addr := net.JoinHostPort(s.host, s.port)
	d := net.Dialer{Timeout: 10 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	if err = conn.SetDeadline(time.Now().Add(30 * time.Second)); err != nil {
		conn.Close()
		return err
	}
	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
			return err
		}
	}

	if err := s.authenticate(c); err != nil {
		return err
	}

	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range recipients {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// buildSMTPMessage renders an rfc 822 message from a postmark email
func buildSMTPMessage(req *postmark.Email) ([]byte, error) {
	for _, v := range []string{req.From, req.To, req.Cc, req.Bcc, req.ReplyTo, req.Tag, req.Subject} {
		if strings.ContainsAny(v, "\r\n") {
			return nil, fmt.Errorf("mail: header value must not contain line breaks")
		}
	}

	var b bytes.Buffer

	fmt.Fprintf(&b, "From: %s\r\n", req.From)
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	fmt.Fprintf(&b, "To: %s\r\n", req.To)
	if req.Cc != "" {
		fmt.Fprintf(&b, "Cc: %s\r\n", req.Cc)
	}
	if req.ReplyTo != "" {
		fmt.Fprintf(&b, "Reply-To: %s\r\n", req.ReplyTo)
	}
	if req.Tag != "" {
		fmt.Fprintf(&b, "X-Tag: %s\r\n", req.Tag)
	}
	fmt.Fprintf(&b, "Subject: %s\r\n", encodeSubject(req.Subject))
	b.WriteString("MIME-Version: 1.0\r\n")

	switch {
	case req.TextBody != "" && req.HTMLBody != "":
		fmt.Fprintf(&b, "Content-Type: multipart/alternative; boundary=%s\r\n\r\n", boundary)
		fmt.Fprintf(&b, "--%s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s", boundary, req.TextBody)
		fmt.Fprintf(&b, "\r\n--%s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s", boundary, req.HTMLBody)
		fmt.Fprintf(&b, "\r\n--%s--\r\n", boundary)
	case req.TextBody != "":
		b.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		b.WriteString(req.TextBody)
	default:
		b.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
		b.WriteString(req.HTMLBody)
	}

	return b.Bytes(), nil
}

// encodeSubject applies rfc 2047 base64 encoding to non-ascii subjects
func encodeSubject(subject string) string {
	if isASCII(subject) {
		return subject
	}
	return "=?utf-8?b?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
