package mail

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"math/big"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type fakeSMTPConfig struct {
	implicitTLS bool
	startTLS    bool
	rejectRcpt  bool
}

type fakeSMTP struct {
	config   fakeSMTPConfig
	listener net.Listener
	tlsCfg   *tls.Config
	mu       sync.Mutex
	auth     string
	from     string
	rcpts    []string
	data     string
	dot      bool
	upgraded bool
}

func newFakeSMTP(t *testing.T, config fakeSMTPConfig) *fakeSMTP {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	fake := &fakeSMTP{
		config:   config,
		listener: listener,
		tlsCfg:   &tls.Config{Certificates: []tls.Certificate{genCert(t)}},
	}
	if config.implicitTLS {
		fake.listener = tls.NewListener(listener, fake.tlsCfg)
	}
	go fake.serve()

	return fake
}

func (f *fakeSMTP) serve() {
	conn, err := f.listener.Accept()
	if err != nil {
		return
	}
	f.handle(conn)
}

func (f *fakeSMTP) handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	send := func(lines ...string) {
		for _, line := range lines {
			writer.WriteString(line + "\r\n")
		}
		writer.Flush()
	}
	send("220 fake ESMTP")
	inData := false
	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line := strings.TrimRight(raw, "\r\n")
		if inData {
			if line == "." {
				inData = false
				f.mu.Lock()
				f.dot = true
				f.mu.Unlock()
				send("250 message accepted")
				continue
			}
			f.mu.Lock()
			f.data += line + "\r\n"
			f.mu.Unlock()
			continue
		}

		command := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(command, "EHLO"), strings.HasPrefix(command, "HELO"):
			send(f.greeting()...)
		case strings.HasPrefix(command, "STARTTLS"):
			send("220 ready to start TLS")
			tlsConn := tls.Server(conn, f.tlsCfg)
			if err := tlsConn.Handshake(); err != nil {
				return
			}
			f.mu.Lock()
			f.upgraded = true
			f.mu.Unlock()
			conn = tlsConn
			reader = bufio.NewReader(conn)
			writer = bufio.NewWriter(conn)
		case strings.HasPrefix(command, "AUTH"):
			f.mu.Lock()
			f.auth = strings.TrimSpace(line[len("AUTH"):])
			f.mu.Unlock()
			send("235 authenticated")
		case strings.HasPrefix(command, "MAIL FROM:"):
			f.mu.Lock()
			f.from = angleAddr(line)
			f.mu.Unlock()
			send("250 sender ok")
		case strings.HasPrefix(command, "RCPT TO:"):
			if f.config.rejectRcpt {
				send("550 recipient rejected")
				continue
			}
			f.mu.Lock()
			f.rcpts = append(f.rcpts, angleAddr(line))
			f.mu.Unlock()
			send("250 recipient ok")
		case strings.HasPrefix(command, "DATA"):
			inData = true
			send("354 end with CRLF.CRLF")
		case strings.HasPrefix(command, "QUIT"):
			send("221 bye")
			return
		default:
			send("250 ok")
		}
	}
}

func (f *fakeSMTP) greeting() []string {
	lines := []string{"250-fake ESMTP"}
	if f.config.startTLS {
		lines = append(lines, "250-STARTTLS")
	}

	return append(lines, "250 AUTH PLAIN")
}

func (f *fakeSMTP) port() string {
	_, port, err := net.SplitHostPort(f.listener.Addr().String())
	if err != nil {
		return ""
	}

	return port
}

func (f *fakeSMTP) close() {
	f.listener.Close()
}

func (f *fakeSMTP) authPayload() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.auth
}

func (f *fakeSMTP) mailFrom() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.from
}

func (f *fakeSMTP) recipients() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.rcpts
}

func (f *fakeSMTP) message() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.data
}

func (f *fakeSMTP) dotSeen() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.dot
}

func (f *fakeSMTP) isUpgraded() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.upgraded
}

func angleAddr(line string) string {
	start := strings.Index(line, "<")
	end := strings.Index(line, ">")
	if start < 0 || end <= start {
		return ""
	}

	return line[start+1 : end]
}

func genCert(t *testing.T) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), nil)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

func insecureTLS(t *testing.T, sender Sender) *SMTPSender {
	t.Helper()
	smtpSender, ok := sender.(*SMTPSender)
	if !ok {
		t.Fatal("sender is not an SMTPSender")
	}
	smtpSender.tlsCfg.InsecureSkipVerify = true

	return smtpSender
}

type smtpSuite struct {
	suite.Suite
}

func (s *smtpSuite) TestSend() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{})
	defer fake.close()
	sender := NewSMTP("127.0.0.1", fake.port(), "login", "password")
	req := s.req()
	req.Cc = "cc@example.com"
	req.Bcc = "bcc@example.com"
	req.Tag = "confirmation"

	s.NoError(sender.Send(s.T().Context(), req))
	payload := base64.StdEncoding.EncodeToString([]byte("\x00login\x00password"))
	s.Equal([]string{"PLAIN", payload}, strings.Fields(fake.authPayload()))
	s.Equal("from@example.com", fake.mailFrom())
	s.Equal([]string{"to@example.com", "cc@example.com", "bcc@example.com"}, fake.recipients())
	s.True(fake.dotSeen())
	msg := fake.message()
	s.Contains(msg, "Subject: subject\r\n")
	s.Contains(msg, "X-Tag: confirmation\r\n")
	s.Contains(msg, "text body")
	s.NotContains(msg, "Bcc")
}

func (s *smtpSuite) TestNewSMTP() {
	implicit, ok := NewSMTP("127.0.0.1", smtpImplicitTLSPort, "", "").(*SMTPSender)
	s.Require().True(ok)
	s.True(implicit.implicitTLS)

	starttls, ok := NewSMTP("127.0.0.1", "587", "", "").(*SMTPSender)
	s.Require().True(ok)
	s.False(starttls.implicitTLS)
}

func (s *smtpSuite) TestSend_DisplayNameFrom() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{})
	defer fake.close()
	sender := NewSMTP("127.0.0.1", fake.port(), "", "")
	req := s.req()
	req.From = "Buscarron <bot@example.com>"

	s.NoError(sender.Send(s.T().Context(), req))
	s.Equal("bot@example.com", fake.mailFrom())
	s.Contains(fake.message(), "From: Buscarron <bot@example.com>\r\n")
}

func (s *smtpSuite) TestSend_NoAuth() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{})
	defer fake.close()
	sender := NewSMTP("127.0.0.1", fake.port(), "", "")

	s.NoError(sender.Send(s.T().Context(), s.req()))
	s.Empty(fake.authPayload())
	s.Equal("from@example.com", fake.mailFrom())
}

func (s *smtpSuite) TestSend_ImplicitTLS() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{implicitTLS: true, startTLS: true})
	defer fake.close()
	sender := insecureTLS(s.T(), NewSMTP("127.0.0.1", fake.port(), "login", "password"))
	sender.implicitTLS = true

	s.NoError(sender.Send(s.T().Context(), s.req()))
	s.False(fake.isUpgraded())
	s.Equal("from@example.com", fake.mailFrom())
	s.Contains(fake.message(), "text body")
}

func (s *smtpSuite) TestSend_StartTLS() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{startTLS: true})
	defer fake.close()
	sender := insecureTLS(s.T(), NewSMTP("127.0.0.1", fake.port(), "login", "password"))

	s.NoError(sender.Send(s.T().Context(), s.req()))
	s.True(fake.isUpgraded())
	s.Equal("from@example.com", fake.mailFrom())
	s.Contains(fake.message(), "text body")
}

func (s *smtpSuite) TestSend_CanceledContext() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{})
	defer fake.close()
	sender := NewSMTP("127.0.0.1", fake.port(), "", "")
	ctx, cancel := context.WithCancel(s.T().Context())
	cancel()

	s.Error(sender.Send(ctx, s.req()))
	s.Empty(fake.mailFrom())
}

func (s *smtpSuite) TestSend_RcptRejected() {
	fake := newFakeSMTP(s.T(), fakeSMTPConfig{rejectRcpt: true})
	defer fake.close()
	sender := NewSMTP("127.0.0.1", fake.port(), "", "")

	s.ErrorContains(sender.Send(s.T().Context(), s.req()), "550")
	s.NotContains(fake.message(), "text body")
}

func (s *smtpSuite) req() *postmark.Email {
	return &postmark.Email{From: "from@example.com", To: "to@example.com", Subject: "subject", TextBody: "text body"}
}

func TestSMTP(t *testing.T) {
	suite.Run(t, new(smtpSuite))
}
