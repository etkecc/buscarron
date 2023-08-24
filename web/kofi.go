package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/id"
)

type SenderByEmail interface {
	SendByEmail(roomID id.RoomID, email, message string, reactions ...string) bool
}

type KoFiConfig struct {
	VerificationToken string
	Logger            *zerolog.Logger
	Sender            SenderByEmail
	Rooms             []id.RoomID
}

type kofi struct {
	token  string
	log    *zerolog.Logger
	sender SenderByEmail
	rooms  []id.RoomID
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

func NewKoFiHandler(cfg *KoFiConfig) *kofi {
	return &kofi{
		token:  cfg.VerificationToken,
		log:    cfg.Logger,
		sender: cfg.Sender,
		rooms:  cfg.Rooms,
	}
}

func (k *kofi) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			k.log.Error().Err(err).Msg("cannot parse form of a ko-fi request")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		raw := r.PostFormValue("data")
		var data *kofiRequest
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			k.log.Error().Err(err).Msg("cannot parse json data of a ko-fi request")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if data.VerificationToken != k.token {
			k.log.Error().Str("provided_token", data.VerificationToken).Msg("verification token is invalid")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		var txt strings.Builder
		txt.WriteString(data.Type)
		txt.WriteString(" payment received!\n\n")

		if data.TierName != nil {
			txt.WriteString("* Tier: ")
			txt.WriteString(*data.TierName)
			txt.WriteString("\n")
		}
		txt.WriteString("* Amount: ")
		txt.WriteString(data.Amount)
		txt.WriteString(" ")
		txt.WriteString(data.Currency)
		txt.WriteString("\n")

		txt.WriteString("* Transaction ID: ")
		txt.WriteString(data.KofiTransactionID)
		txt.WriteString("\n\n")

		if data.Message != nil {
			txt.WriteString("> ")
			txt.WriteString(*data.Message)
			txt.WriteString("\n> --")
			txt.WriteString(data.FromName)
		}

		message := txt.String()
		k.log.Debug().Msg(message)
		for _, roomID := range k.rooms {
			if ok := k.sender.SendByEmail(roomID, data.Email, message, "💸"); ok {
				k.log.Info().Str("roomID", roomID.String()).Msg("successfully sent ko-fi update into the room by email")
				return
			}
		}
		w.Write([]byte("ok")) //nolint:errcheck
	}
}
