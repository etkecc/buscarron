package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/go/psd"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type senderByEmail interface {
	Send(ctx context.Context, roomID id.RoomID, message string, attributes map[string]interface{}) id.EventID
	SendByEmail(ctx context.Context, roomID id.RoomID, email, message string, reactions ...string) map[string]any
	FindEventBy(ctx context.Context, roomID id.RoomID, field, value string, fromToken ...string) *event.Event
}

type KoFiConfig struct {
	VerificationToken string
	Room              id.RoomID
	Sender            senderByEmail
	PaidMarker        func(context.Context, string, string, string)
	Rooms             []id.RoomID
	PSD               *psd.Client
}

type kofi struct {
	token        string
	psdc         *psd.Client
	sender       senderByEmail
	markPaid     func(context.Context, string, string, string)
	rooms        []id.RoomID
	fallbackRoom id.RoomID
}

type kofiRequest struct {
	VerificationToken          string    `json:"verification_token"`
	MessageID                  string    `json:"message_id"`
	Timestamp                  time.Time `json:"timestamp"`
	Type                       string    `json:"type"`
	IsPublic                   bool      `json:"is_public"`
	FromName                   string    `json:"from_name"`
	Message                    *string   `json:"message"`
	Amount                     string    `json:"amount"`
	URL                        string    `json:"url"`
	Email                      string    `json:"email"`
	Currency                   string    `json:"currency"`
	IsSubscriptionPayment      bool      `json:"is_subscription_payment"`
	IsFirstSubscriptionPayment bool      `json:"is_first_subscription_payment"`
	KofiTransactionID          string    `json:"kofi_transaction_id"`
	TierName                   *string   `json:"tier_name"`
}

func (r *kofiRequest) getStatus(ctx context.Context, psdc *psd.Client) (bool, string) {
	targets, err := psdc.GetWithContext(ctx, r.Email)
	if err != nil {
		return false, ""
	}
	if len(targets) == 0 {
		return false, ""
	}
	return true, targets[0].GetDomain()
}

func (r *kofiRequest) Text(ctx context.Context, psdc *psd.Client) string {
	var txt strings.Builder
	txt.WriteString(r.Type)
	txt.WriteString(" payment received!\n\n")

	ok, domain := r.getStatus(ctx, psdc)
	txt.WriteString("* Email: ")
	if ok {
		txt.WriteString("👤")
	}
	txt.WriteString(r.Email)
	txt.WriteString("\n")

	if ok {
		txt.WriteString("* Host: 👥")
		txt.WriteString(domain)
		txt.WriteString("\n")
	}

	if r.TierName != nil {
		txt.WriteString("* Tier: ")
		txt.WriteString(*r.TierName)
		txt.WriteString("\n")
	}
	txt.WriteString("* Amount: ")
	txt.WriteString(r.Amount)
	txt.WriteString(" ")
	txt.WriteString(r.Currency)
	txt.WriteString("\n")

	txt.WriteString("* Transaction ID: ")
	txt.WriteString(r.KofiTransactionID)
	txt.WriteString("\n\n")

	if r.Message != nil {
		txt.WriteString("> ")
		txt.WriteString(*r.Message)
		txt.WriteString("\n")
	}
	txt.WriteString("> --")
	txt.WriteString(r.FromName)

	return txt.String()
}

func (r *kofiRequest) Logger(ctx context.Context) *zerolog.Logger {
	ctxlog := zerolog.Ctx(ctx).With().
		Str("email", r.Email).
		Str("type", r.Type).
		Bool("is_subscription", r.IsSubscriptionPayment).
		Bool("is_first", r.IsFirstSubscriptionPayment).
		Logger()
	return &ctxlog
}

func NewKoFi(cfg *KoFiConfig) *kofi {
	return &kofi{
		token:        cfg.VerificationToken,
		psdc:         cfg.PSD,
		sender:       cfg.Sender,
		markPaid:     cfg.PaidMarker,
		rooms:        cfg.Rooms,
		fallbackRoom: cfg.Room,
	}
}

func (k *kofi) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		log := zerolog.Ctx(c.Request().Context())
		raw := c.FormValue("data")
		var data *kofiRequest
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			log.Error().Err(err).Msg("cannot parse json data of a ko-fi request")
			return c.NoContent(http.StatusBadRequest)
		}

		if data.VerificationToken != k.token {
			log.Error().Str("provided_token", data.VerificationToken).Msg("verification token is invalid")
			return c.NoContent(http.StatusUnauthorized)
		}

		return k.send(c, data)
	}
}

func (k *kofi) send(c echo.Context, data *kofiRequest) error {
	ctx := c.Request().Context()
	log := data.Logger(ctx)
	message := data.Text(ctx, k.psdc)
	// if one-off - just send the message
	if data.Type != "Subscription" {
		k.fallback(ctx, data, message)
		return c.NoContent(http.StatusOK)
	}

	// if subscription - send the message only if it's the first payment
	if !data.IsFirstSubscriptionPayment {
		return c.NoContent(http.StatusOK)
	}

	for _, roomID := range k.rooms {
		if raw := k.sender.SendByEmail(ctx, roomID, data.Email, message, "💸"); raw != nil {
			log.Info().Str("roomID", roomID.String()).Msg("successfully sent ko-fi update into the room by email")
			domain, ok := raw["domain"].(string)
			baseDomain, _ := raw["base_domain"].(string)
			if ok && k.markPaid != nil {
				k.markPaid(ctx, domain, baseDomain, data.Amount)
			} else {
				log.Error().Any("domain", domain).Msg("cannot mark as paid, domain is not a string")
			}
			return c.NoContent(http.StatusOK)
		}
	}
	k.fallback(ctx, data, message)
	return c.NoContent(http.StatusOK)
}

func (k *kofi) fallback(ctx context.Context, data *kofiRequest, message string) {
	if k.sender.FindEventBy(ctx, k.fallbackRoom, "kofi_id", data.KofiTransactionID) != nil {
		return
	}
	k.sender.Send(ctx, k.fallbackRoom, message, map[string]interface{}{"kofi_id": data.KofiTransactionID})
}
