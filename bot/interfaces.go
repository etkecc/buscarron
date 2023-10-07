package bot

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Linkpearl interface
type Linkpearl interface {
	Send(roomID id.RoomID, content interface{}) (id.EventID, error)
	SendFile(roomID id.RoomID, req *mautrix.ReqUploadMedia, msgtype event.MessageType, relations ...*event.RelatesTo) error
	SendNotice(roomID id.RoomID, message string, relates ...*event.RelatesTo)
	FindEventBy(roomID id.RoomID, field, value string, fromToken ...string) *event.Event
	Start(optionalStatusMsg ...string) error
	GetClient() *mautrix.Client
	Stop()
}
