package bot

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Linkpearl interface
type Linkpearl interface {
	Send(roomID id.RoomID, content interface{}) (id.EventID, error)
	SendFile(roomID id.RoomID, req *mautrix.ReqUploadMedia, msgtype event.MessageType, relation *event.RelatesTo) error
	Start(optionalStatusMsg ...string) error
	Stop()
}
