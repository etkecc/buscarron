package mail

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"mime"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mattevans/postmark-go"
)

// buildMessage renders the email as an RFC 5322 message
func buildMessage(req *postmark.Email) ([]byte, error) {
	if err := guardHeaders(req); err != nil {
		return nil, err
	}
	if req.From == "" {
		return nil, errors.New("mail: empty from address")
	}
	if req.TextBody == "" && req.HTMLBody == "" {
		return nil, errors.New("mail: empty message body")
	}

	var buf bytes.Buffer
	writeHeaders(&buf, req)
	writeBody(&buf, req)

	return buf.Bytes(), nil
}

// guardHeaders rejects values that could inject extra headers or lines
func guardHeaders(req *postmark.Email) error {
	for _, value := range []string{req.From, req.To, req.Cc, req.Bcc, req.ReplyTo, req.Tag, req.Subject} {
		if strings.ContainsAny(value, "\r\n") {
			return errors.New("mail: header value must not contain line breaks")
		}
	}
	for _, header := range req.Headers {
		if strings.ContainsAny(header.Name+header.Value, "\r\n") {
			return errors.New("mail: header value must not contain line breaks")
		}
		if strings.Contains(header.Name, ":") {
			return errors.New("mail: header name must not contain colon")
		}
	}

	return nil
}

// writeHeaders writes the message headers
func writeHeaders(buf *bytes.Buffer, req *postmark.Email) {
	fmt.Fprintf(buf, "From: %s\r\n", req.From)
	fmt.Fprintf(buf, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	if req.To != "" {
		fmt.Fprintf(buf, "To: %s\r\n", req.To)
	}
	if req.Cc != "" {
		fmt.Fprintf(buf, "Cc: %s\r\n", req.Cc)
	}
	if req.ReplyTo != "" {
		fmt.Fprintf(buf, "Reply-To: %s\r\n", req.ReplyTo)
	}
	if req.Tag != "" {
		fmt.Fprintf(buf, "X-Tag: %s\r\n", req.Tag)
	}
	fmt.Fprintf(buf, "Subject: %s\r\n", mime.BEncoding.Encode("utf-8", req.Subject))
	for _, header := range req.Headers {
		fmt.Fprintf(buf, "%s: %s\r\n", header.Name, header.Value)
	}
	if id := messageID(req.From); id != "" {
		fmt.Fprintf(buf, "Message-ID: %s\r\n", id)
	}
	buf.WriteString("MIME-Version: 1.0\r\n")
}

// writeBody writes the body part, multipart/alternative when both bodies are set
func writeBody(buf *bytes.Buffer, req *postmark.Email) {
	text := crlf(req.TextBody)
	html := crlf(req.HTMLBody)
	switch {
	case text != "" && html != "":
		boundary := rand.Text()
		fmt.Fprintf(buf, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)
		fmt.Fprintf(buf, "--%s\r\nContent-Type: text/plain; charset=utf-8\r\n%s\r\n%s\r\n", boundary, cte(text), text)
		fmt.Fprintf(buf, "--%s\r\nContent-Type: text/html; charset=utf-8\r\n%s\r\n%s\r\n", boundary, cte(html), html)
		fmt.Fprintf(buf, "--%s--\r\n", boundary)
	case text != "":
		fmt.Fprintf(buf, "Content-Type: text/plain; charset=utf-8\r\n%s\r\n%s", cte(text), text)
	default:
		fmt.Fprintf(buf, "Content-Type: text/html; charset=utf-8\r\n%s\r\n%s", cte(html), html)
	}
}

// cte returns the 8bit encoding header when the body has non-ASCII bytes
func cte(body string) string {
	if isASCII(body) {
		return ""
	}

	return "Content-Transfer-Encoding: 8bit\r\n"
}

// isASCII reports whether the body has only 7-bit bytes
func isASCII(body string) bool {
	for i := range len(body) {
		if body[i] >= utf8.RuneSelf {
			return false
		}
	}

	return true
}

// messageID builds a Message-ID from the sender domain
func messageID(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return ""
	}
	_, domain, ok := strings.Cut(addr.Address, "@")
	if !ok {
		return ""
	}

	return "<" + rand.Text() + "@" + domain + ">"
}

// crlf normalizes line breaks to CRLF
func crlf(body string) string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\r", "\n")

	return strings.ReplaceAll(body, "\n", "\r\n")
}
