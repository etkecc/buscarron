package bot

import (
	"database/sql"

	"gitlab.com/etke.cc/linkpearl/store"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Linkpearl interface
type Linkpearl interface {
	GetClient() *mautrix.Client
	GetDB() *sql.DB
	GetMachine() *crypto.OlmMachine
	GetStore() *store.Store
	OnEvent(callback mautrix.EventHandler)
	OnEventType(eventType event.Type, callback mautrix.EventHandler)
	OnSync(callback mautrix.SyncHandler)
	Send(roomID id.RoomID, content interface{}) (id.EventID, error)
	Start() error
	Stop()
}
