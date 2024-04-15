package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/buscarron/utils"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

// Bot represents matrix bot
type Bot struct {
	mu sync.Mutex
	lp Linkpearl
}

// New creates a new matrix bot
func New(lp Linkpearl) *Bot {
	return &Bot{
		lp: lp,
	}
}

// Error message to the log and matrix room
func (b *Bot) Error(ctx context.Context, roomID id.RoomID, message string, args ...any) {
	zerolog.Ctx(ctx).Error().Msgf(message, args...)

	if b.lp == nil {
		return
	}
	b.Send(ctx, roomID, "ERROR: "+fmt.Sprintf(message, args...), nil)
}

// Send message to the room
//
//nolint:unparam // return value is used, but called from interfaces
func (b *Bot) Send(ctx context.Context, roomID id.RoomID, message string, attributes map[string]any) id.EventID {
	span := utils.StartSpan(ctx, "bot.Send")
	defer span.Finish()

	parsed := format.RenderMarkdown(message, true, true)
	parsed.MsgType = event.MsgNotice
	content := event.Content{
		Raw:    attributes,
		Parsed: &parsed,
	}

	b.mu.Lock()
	eventID, err := b.lp.Send(span.Context(), roomID, &content)
	b.mu.Unlock()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("roomID", roomID.String()).Msg("cannot send message")
	}
	return eventID
}

// SendFile for the room
func (b *Bot) SendFile(ctx context.Context, roomID id.RoomID, file *mautrix.ReqUploadMedia, relations ...*event.RelatesTo) {
	span := utils.StartSpan(ctx, "linkpearl.SendFile")
	defer span.Finish()

	b.mu.Lock()
	err := b.lp.SendFile(span.Context(), roomID, file, event.MsgFile, relations...)
	b.mu.Unlock()
	if err != nil {
		b.Error(span.Context(), roomID, "cannot upload file: %v", err)
		return
	}
}

// Start performs matrix /sync
func (b *Bot) Start() {
	ctx := utils.NewContext()
	if err := b.lp.Start(ctx); err != nil {
		zerolog.Ctx(ctx).Panic().Err(err).Msg("matrix bot crashed")
	}
}

// Stop the bot
func (b *Bot) Stop() {
	b.lp.Stop(utils.NewContext())
}
