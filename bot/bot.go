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

// SendByEmail sends the message into the room as thread reply by email
func (b *Bot) SendByEmail(ctx context.Context, roomID id.RoomID, email, message string, reactions ...string) map[string]any {
	log := zerolog.Ctx(ctx)
	span := utils.StartSpan(ctx, "bot.SendByEmail")
	defer span.Finish()

	evt := b.FindEventBy(span.Context(), roomID, "email", email)
	if evt == nil {
		log.Warn().Str("roomID", roomID.String()).Msg("event by email was not found in that room")
		return nil
	}

	content := format.RenderMarkdown(message, true, true)
	content.MsgType = event.MsgNotice
	content.SetRelatesTo(&event.RelatesTo{
		Type:    event.RelThread,
		EventID: evt.ID,
	})

	b.mu.Lock()
	sendSpan := utils.StartSpan(span.Context(), "linkpearl.Send")
	_, err := b.lp.Send(sendSpan.Context(), roomID, content)
	sendSpan.Finish()
	b.mu.Unlock()
	if err != nil {
		log.Warn().Err(err).Str("roomID", roomID.String()).Str("threadID", evt.ID.String()).Msg("cannot send a message by email")
		return nil
	}

	if len(reactions) > 0 {
		for _, reaction := range reactions {
			reactionSpan := utils.StartSpan(span.Context(), "mautrix.SendReaction")
			_, err = b.lp.GetClient().SendReaction(reactionSpan.Context(), roomID, evt.ID, reaction)
			reactionSpan.Finish()
			if err != nil {
				log.Warn().Err(err).Str("roomID", roomID.String()).Str("eventID", evt.ID.String()).Msg("cannot send a reaction")
			}
		}
	}

	return evt.Content.Raw
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

// FindEventBy is wrapper around lp.FindEventBy
func (b *Bot) FindEventBy(ctx context.Context, roomID id.RoomID, field, value string, fromToken ...string) *event.Event {
	span := utils.StartSpan(ctx, "linkpearl.FindEventBy")
	defer span.Finish()
	return b.lp.FindEventBy(span.Context(), roomID, field, value, fromToken...)
}

// Start performs matrix /sync
func (b *Bot) Start() {
	if err := b.lp.Start(); err != nil {
		zerolog.Ctx(utils.NewContext()).Panic().Err(err).Msg("matrix bot crashed")
	}
}

// Stop the bot
func (b *Bot) Stop() {
	b.lp.Stop()
}
