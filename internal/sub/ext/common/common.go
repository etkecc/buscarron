package common

import (
	"bytes"
	"context"
	"net"
	"text/template"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Validator interface
type Validator interface {
	A(string) bool
	MX(string) bool
	NS(string, ...string) bool
	CNAME(string) bool
	Email(string, string, ...net.IP) bool
	Domain(string) bool
	DomainString(string) bool
	GetBase(string) string
}

// Sender interface to send messages
type Sender interface {
	Send(context.Context, id.RoomID, string, map[string]any) id.EventID
	SendFile(context.Context, id.RoomID, *mautrix.ReqUploadMedia, ...*event.RelatesTo)
}

func ParseTemplate(tplString string, data map[string]string) (string, error) {
	var result bytes.Buffer
	tpl, err := template.New("template").Parse(tplString)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
