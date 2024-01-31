package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type senderByEmail interface {
	Send(roomID id.RoomID, message string, attributes map[string]interface{}) id.EventID
	SendByEmail(roomID id.RoomID, email, message string, reactions ...string) map[string]any
	FindEventBy(roomID id.RoomID, field, value string, fromToken ...string) *event.Event
}

type KoFiConfig struct {
	VerificationToken string
	Room              id.RoomID
	Logger            *zerolog.Logger
	Sender            senderByEmail
	PaidMarker        func(*zerolog.Logger, string, string, string)
	Rooms             []id.RoomID
}

type kofi struct {
	token        string
	log          *zerolog.Logger
	sender       senderByEmail
	markPaid     func(*zerolog.Logger, string, string, string)
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

func (r *kofiRequest) Text() string {
	var txt strings.Builder
	txt.WriteString(r.Type)
	txt.WriteString(" payment received!\n\n")

	txt.WriteString("* Email: ")
	txt.WriteString(r.Email)
	txt.WriteString("\n")

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

func (r *kofiRequest) Logger(log *zerolog.Logger) *zerolog.Logger {
	ctxlog := log.With().
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
		log:          cfg.Logger,
		sender:       cfg.Sender,
		markPaid:     cfg.PaidMarker,
		rooms:        cfg.Rooms,
		fallbackRoom: cfg.Room,
	}
}

func (k *kofi) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		raw := c.FormValue("data")
		var data *kofiRequest
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			k.log.Error().Err(err).Msg("cannot parse json data of a ko-fi request")
			return c.NoContent(http.StatusBadRequest)
		}
		log := data.Logger(k.log)

		if data.VerificationToken != k.token {
			log.Error().Str("provided_token", data.VerificationToken).Msg("verification token is invalid")
			return c.NoContent(http.StatusUnauthorized)
		}

		return k.send(c, data)
	}
}

func (k *kofi) send(c echo.Context, data *kofiRequest) error {
	log := data.Logger(k.log)
	message := data.Text()
	// if one-off - just send the message
	if data.Type != "Subscription" {
		k.fallback(data, message)
		return c.NoContent(http.StatusOK)
	}

	// if subscription - send the message only if it's the first payment
	if !data.IsFirstSubscriptionPayment {
		return c.NoContent(http.StatusOK)
	}

	for _, roomID := range k.rooms {
		if raw := k.sender.SendByEmail(roomID, data.Email, message, "💸"); raw != nil {
			log.Info().Str("roomID", roomID.String()).Msg("successfully sent ko-fi update into the room by email")
			domain, ok := raw["domain"].(string)
			baseDomain, _ := raw["base_domain"].(string)
			if ok && k.markPaid != nil {
				k.markPaid(log, domain, baseDomain, data.Amount)
			} else {
				log.Error().Any("domain", domain).Msg("cannot mark as paid, domain is not a string")
			}
			return c.NoContent(http.StatusOK)
		}
	}
	k.fallback(data, message)
	return c.NoContent(http.StatusOK)
}

func (k *kofi) fallback(data *kofiRequest, message string) {
	if k.sender.FindEventBy(k.fallbackRoom, "kofi_id", data.KofiTransactionID) != nil {
		return
	}
	k.sender.Send(k.fallbackRoom, message, map[string]interface{}{"kofi_id": data.KofiTransactionID})
}
