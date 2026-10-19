package mail

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/mattevans/postmark-go"
)

const (
	// smtpImplicitTLSPort is the port using implicit TLS
	smtpImplicitTLSPort = "465"
	// smtpDialTimeout is the timeout of the TCP dial
	smtpDialTimeout = 10 * time.Second
	// smtpSessionTimeout bounds the whole SMTP session
	smtpSessionTimeout = 30 * time.Second
)

// SMTPSender sends mail via SMTP
type SMTPSender struct {
	host        string
	port        string
	login       string
	password    string
	implicitTLS bool
	tlsCfg      *tls.Config
}

// NewSMTP creates the SMTP transport
func NewSMTP(host, port, login, password string) Sender {
	return &SMTPSender{
		host:        host,
		port:        port,
		login:       login,
		password:    password,
		implicitTLS: port == smtpImplicitTLSPort,
		tlsCfg:      &tls.Config{ServerName: host},
	}
}

// Send sends the email via SMTP
func (s *SMTPSender) Send(ctx context.Context, req *postmark.Email) error {
	recipients, err := parseRecipients(req)
	if err != nil {
		return err
	}
	msg, err := buildMessage(req)
	if err != nil {
		return err
	}

	return s.deliver(ctx, req.From, recipients, msg)
}

// parseRecipients lists the envelope recipients from To, Cc and Bcc
func parseRecipients(req *postmark.Email) ([]string, error) {
	list := make([]string, 0, 3)
	for _, value := range []string{req.To, req.Cc, req.Bcc} {
		if value != "" {
			list = append(list, value)
		}
	}
	if len(list) == 0 {
		return nil, errors.New("mail: no recipients")
	}
	addrs, err := mail.ParseAddressList(strings.Join(list, ", "))
	if err != nil {
		return nil, fmt.Errorf("mail: cannot parse recipients: %w", err)
	}
	if len(addrs) == 0 {
		return nil, errors.New("mail: no recipients")
	}
	recipients := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		recipients = append(recipients, addr.Address)
	}

	return recipients, nil
}

// deliver runs one SMTP session: dial, upgrade, auth, envelope, data
func (s *SMTPSender) deliver(ctx context.Context, from string, recipients []string, msg []byte) error {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return fmt.Errorf("mail: cannot parse sender: %w", err)
	}
	c, err := s.dial(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	if err := s.startTLS(c); err != nil {
		return err
	}
	if err := s.auth(c); err != nil {
		return err
	}
	if err := c.Mail(addr.Address); err != nil {
		return err
	}
	for _, rcpt := range recipients {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}

	return s.sendData(c, msg)
}

// dial opens and greets the SMTP server
func (s *SMTPSender) dial(ctx context.Context) (*smtp.Client, error) {
	conn, err := s.dialConn(ctx)
	if err != nil {
		return nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(smtpSessionTimeout)); err != nil {
		conn.Close()
		return nil, err
	}
	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return c, nil
}

// dialConn opens a TCP connection, TLS-wrapped on the implicit TLS port
func (s *SMTPSender) dialConn(ctx context.Context) (net.Conn, error) {
	addr := net.JoinHostPort(s.host, s.port)
	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	if s.implicitTLS {
		tlsDialer := tls.Dialer{NetDialer: dialer, Config: s.tlsCfg}
		return tlsDialer.DialContext(ctx, "tcp", addr)
	}

	return dialer.DialContext(ctx, "tcp", addr)
}

// startTLS upgrades the connection when the server advertises STARTTLS
func (s *SMTPSender) startTLS(c *smtp.Client) error {
	if s.implicitTLS {
		return nil
	}
	if ok, _ := c.Extension("STARTTLS"); !ok {
		return nil
	}

	return c.StartTLS(s.tlsCfg)
}

// auth authenticates the session when a login is configured
func (s *SMTPSender) auth(c *smtp.Client) error {
	if s.login == "" {
		return nil
	}
	if ok, _ := c.Extension("AUTH"); !ok {
		return errors.New("mail: smtp server does not advertise AUTH")
	}

	return c.Auth(smtp.PlainAuth("", s.login, s.password, s.host))
}

// sendData writes the message and quits the session
func (s *SMTPSender) sendData(c *smtp.Client, msg []byte) error {
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
