package bot

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/etkecc/buscarron/internal/utils"
	"github.com/rs/zerolog"
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

// Enabled returns true if the bot is enabled
func (b *Bot) Enabled() bool {
	if b.lp == nil {
		return false
	}
	if reflect.ValueOf(b.lp).Kind() == reflect.Pointer && reflect.ValueOf(b.lp).IsNil() {
		return false
	}
	return true
}

// Error message to the log and matrix room
func (b *Bot) Error(ctx context.Context, roomID id.RoomID, message string, args ...any) {
	log := zerolog.Ctx(ctx)
	log.Error().Msgf(message, args...)

	if !b.Enabled() {
		log.Warn().Msg("bot is disabled, cannot send error message to the room")
		return
	}

	parsed := format.RenderMarkdown("ERROR: "+fmt.Sprintf(message, args...), true, true)
	parsed.MsgType = event.MsgNotice
	content := event.Content{
		Parsed: &parsed,
	}

	_, err := b.lp.Send(ctx, roomID, &content)
	if err != nil {
		log.Error().Err(err).Str("roomID", roomID.String()).Msg("cannot send error message")
	}
}

// Send message to the room
//
//nolint:unparam // return value is used, but called from interfaces
func (b *Bot) Send(ctx context.Context, roomID id.RoomID, message string, attributes map[string]any) id.EventID {
	log := zerolog.Ctx(ctx)
	if !b.Enabled() {
		log.Warn().Msg("bot is disabled, cannot send message to the room")
		return ""
	}

	span := utils.StartSpan(ctx, "bot.Send")
	defer span.Finish()

	parsed := format.RenderMarkdown(message, true, true)
	parsed.MsgType = event.MsgNotice
	content := event.Content{
		Raw:    attributes,
		Parsed: &parsed,
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	eventID, err := b.lp.Send(span.Context(), roomID, &content)
	if err != nil {
		log.Error().Err(err).Str("roomID", roomID.String()).Msg("cannot send message")
	}
	return eventID
}

// SendFile for the room
func (b *Bot) SendFile(ctx context.Context, roomID id.RoomID, file *mautrix.ReqUploadMedia, relations ...*event.RelatesTo) {
	log := zerolog.Ctx(ctx)
	if !b.Enabled() {
		log.Warn().Msg("bot is disabled, cannot send file to the room")
		return
	}

	span := utils.StartSpan(ctx, "linkpearl.SendFile")
	defer span.Finish()

	b.mu.Lock()
	defer b.mu.Unlock()
	err := b.lp.SendFile(span.Context(), roomID, file, event.MsgFile, relations...)
	if err != nil {
		b.Error(span.Context(), roomID, "cannot upload file: %v", err)
		return
	}
}

// Start performs matrix /sync
func (b *Bot) Start() error {
	ctx := utils.NewContext()
	log := zerolog.Ctx(ctx)
	if !b.Enabled() {
		log.Warn().Msg("bot is disabled, cannot start it")
		return nil
	}

	return b.lp.Start(ctx)
}

// Stop the bot
func (b *Bot) Stop() {
	ctx := utils.NewContext()
	log := zerolog.Ctx(ctx)
	if !b.Enabled() {
		log.Warn().Msg("bot is disabled, cannot stop it")
		return
	}

	b.lp.Stop(ctx)
}
