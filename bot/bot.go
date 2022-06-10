package bot

import (
	"fmt"
	"sync"

	"gitlab.com/etke.cc/go/logger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

// Bot represents matrix bot
type Bot struct {
	sync.Mutex
	log *logger.Logger
	lp  Linkpearl
}

// New creates a new matrix bot
func New(lp Linkpearl, log *logger.Logger) *Bot {
	return &Bot{
		lp:  lp,
		log: log,
	}
}

// Error message to the log and matrix room
func (b *Bot) Error(roomID id.RoomID, message string, args ...interface{}) {
	b.log.Error(message, args...)

	if b.lp == nil {
		return
	}
	b.Lock()
	defer b.Unlock()
	// nolint // if something goes wrong here nobody can help...
	b.lp.Send(roomID, &event.MessageEventContent{
		MsgType: event.MsgNotice,
		Body:    "ERROR: " + fmt.Sprintf(message, args...),
	})
}

// Send message to the room
func (b *Bot) Send(roomID id.RoomID, message string) {
	content := format.RenderMarkdown(message, true, true)
	b.Lock()
	_, err := b.lp.Send(roomID, &content)
	b.Unlock()
	if err != nil {
		b.Error(roomID, "cannot send message: %v", err)
	}
}

// SendFile for the room
func (b *Bot) SendFile(roomID id.RoomID, file *mautrix.ReqUploadMedia) {
	b.Lock()
	resp, err := b.lp.GetClient().UploadMedia(*file)
	b.Unlock()
	if err != nil {
		b.Error(roomID, "cannot upload file: %v", err)
		return
	}

	b.Lock()
	_, err = b.lp.Send(roomID, &event.MessageEventContent{
		MsgType: event.MsgFile,
		Body:    file.FileName,
		URL:     resp.ContentURI.CUString(),
	})
	b.Unlock()
	if err != nil {
		b.Error(roomID, "cannot send a message with uploaded file: %v", err)
	}
}

// Start performs matrix /sync
func (b *Bot) Start() {
	if err := b.lp.Start(); err != nil {
		b.log.Fatal("matrix bot crashed: %v", err)
		return
	}
}

// Stop the bot
func (b *Bot) Stop() {
	b.lp.Stop()
}
