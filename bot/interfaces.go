package bot

import (
	"context"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Linkpearl interface
type Linkpearl interface {
	Send(ctx context.Context, roomID id.RoomID, content any) (id.EventID, error)
	SendFile(ctx context.Context, roomID id.RoomID, req *mautrix.ReqUploadMedia, msgtype event.MessageType, relations ...*event.RelatesTo) error
	SendNotice(ctx context.Context, roomID id.RoomID, message string, relates ...*event.RelatesTo)
	FindEventBy(ctx context.Context, roomID id.RoomID, field, value string, fromToken ...string) *event.Event
	Start(optionalStatusMsg ...string) error
	GetClient() *mautrix.Client
	Stop()
}
