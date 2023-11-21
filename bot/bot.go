package bot

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

// Bot represents matrix bot
type Bot struct {
	sync.Mutex
	log *zerolog.Logger
	lp  Linkpearl
}

// New creates a new matrix bot
func New(lp Linkpearl, log *zerolog.Logger) *Bot {
	return &Bot{
		lp:  lp,
		log: log,
	}
}

// Error message to the log and matrix room
func (b *Bot) Error(roomID id.RoomID, message string, args ...interface{}) {
	b.log.Error().Msgf(message, args...)

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
func (b *Bot) Send(roomID id.RoomID, message string, attributes map[string]interface{}) id.EventID {
	parsed := format.RenderMarkdown(message, true, true)
	content := event.Content{
		Raw:    attributes,
		Parsed: &parsed,
	}

	b.Lock()
	eventID, err := b.lp.Send(roomID, &content)
	b.Unlock()
	if err != nil {
		b.log.Error().Err(err).Str("roomID", roomID.String()).Msg("cannot send message")
	}
	return eventID
}

// SendByEmail sends the message into the room as thread reply by email
func (b *Bot) SendByEmail(roomID id.RoomID, email string, message string, reactions ...string) bool {
	evt := b.lp.FindEventBy(roomID, "email", email)
	if evt == nil {
		b.log.Warn().Str("roomID", roomID.String()).Msg("event by email was not found in that room")
		return false
	}

	content := format.RenderMarkdown(message, true, true)
	content.SetRelatesTo(&event.RelatesTo{
		Type:    event.RelThread,
		EventID: evt.ID,
	})

	b.Lock()
	_, err := b.lp.Send(roomID, content)
	b.Unlock()
	if err != nil {
		b.log.Warn().Err(err).Str("roomID", roomID.String()).Str("threadID", evt.ID.String()).Msg("cannot send a message by email")
		return false
	}

	if len(reactions) > 0 {
		for _, reaction := range reactions {
			_, err = b.lp.GetClient().SendReaction(roomID, evt.ID, reaction)
			if err != nil {
				b.log.Warn().Err(err).Str("roomID", roomID.String()).Str("eventID", evt.ID.String()).Msg("cannot send a reaction")
			}
		}
	}

	return true
}

// SendFile for the room
func (b *Bot) SendFile(roomID id.RoomID, file *mautrix.ReqUploadMedia, relations ...*event.RelatesTo) {
	b.Lock()
	err := b.lp.SendFile(roomID, file, event.MsgFile, relations...)
	b.Unlock()
	if err != nil {
		b.Error(roomID, "cannot upload file: %v", err)
		return
	}
}

// FindEventBy is wrapper around lp.FindEventBy
func (b *Bot) FindEventBy(roomID id.RoomID, field, value string, fromToken ...string) *event.Event {
	return b.lp.FindEventBy(roomID, field, value, fromToken...)
}

// Start performs matrix /sync
func (b *Bot) Start() {
	if err := b.lp.Start(); err != nil {
		b.log.Panic().Err(err).Msg("matrix bot crashed")
		return
	}
}

// Stop the bot
func (b *Bot) Stop() {
	b.lp.Stop()
}
