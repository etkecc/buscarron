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
	SendByEmail(roomID id.RoomID, email, message string, reactions ...string) bool
	FindEventBy(roomID id.RoomID, field, value string, fromToken ...string) *event.Event
}

type KoFiConfig struct {
	VerificationToken string
	Room              id.RoomID
	Logger            *zerolog.Logger
	Sender            senderByEmail
	Rooms             []id.RoomID
}

type kofi struct {
	token        string
	log          *zerolog.Logger
	sender       senderByEmail
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

		// not a first subscription payment = ignore
		if !data.IsSubscriptionPayment || !data.IsFirstSubscriptionPayment {
			log.Info().Msg("not a first subscription payment, ignoring")
			return c.NoContent(http.StatusOK)
		}

		message := data.Text()
		for _, roomID := range k.rooms {
			if ok := k.sender.SendByEmail(roomID, data.Email, message, "💸"); ok {
				log.Info().Str("roomID", roomID.String()).Msg("successfully sent ko-fi update into the room by email")
				return c.NoContent(http.StatusOK)
			}
			k.fallback(data, message)
		}
		return c.NoContent(http.StatusOK)
	}
}

func (k *kofi) fallback(data *kofiRequest, message string) {
	if k.sender.FindEventBy(k.fallbackRoom, "kofi_id", data.KofiTransactionID) != nil {
		return
	}
	k.sender.Send(k.fallbackRoom, message, map[string]interface{}{"kofi_id": data.KofiTransactionID})
}
