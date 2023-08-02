package web

import (
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/buscarron/metrics"
)

type banhandler struct {
	store *lru.Cache[string, struct{}]
	log   *zerolog.Logger
}

var donotban = []string{"/favicon.ico", "/robots.txt", "/metrics", "/_domain"}

// NewBanHandler creates banhandler
func NewBanHandler(size int, banlist []string, log *zerolog.Logger) *banhandler {
	store, err := lru.New[string, struct{}](size)
	if err != nil {
		log.Error().Err(err).Msg("cannot init cache")
	}
	log.Info().Strs("ids", banlist).Msg("preemptive banning...")
	for _, id := range banlist {
		id = strings.TrimSpace(id)
		store.Add(id, struct{}{})
		metrics.BanUser("persistent", "-")
	}

	return &banhandler{store, log}
}

// Ban by request data
func (b *banhandler) Ban(r *http.Request, reason string) {
	for _, exclude := range donotban {
		if r.URL.Path == exclude {
			return
		}
	}
	ctxIDv := r.Context().Value(ctxID)
	id, ok := ctxIDv.(string)
	if !ok {
		b.log.Error().Any("ctxID", ctxIDv).Msg("cannot convert ctxID to string")
		return
	}
	if id == "" || id == "1" {
		b.log.Error().Any("ctxID", ctxIDv).Msg("hashed IP is empty")
		return
	}
	b.store.Add(id, struct{}{})
	ctxNameV := r.Context().Value(ctxName)
	name, ok := ctxNameV.(string)
	if !ok {
		b.log.Error().Any("ctxName", ctxNameV).Msg("cannot convert ctxName to string")
	}
	metrics.BanUser(reason, name)

	b.log.Info().Str("method", r.Method).Str("path", r.URL.Path).Str("id", id).Msg("has been banned")
}

// Handle to check for ban
func (b *banhandler) Handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxIDv := r.Context().Value(ctxID)
		id, ok := ctxIDv.(string)
		if !ok {
			b.log.Error().Any("ctxID", ctxIDv).Msg("cannot convert ctxID to string")
		}
		id = strings.TrimSpace(id)
		if b.store.Contains(id) {
			metrics.BanRequest()
			b.log.Debug().Str("method", r.Method).Str("url", r.URL.String()).Str("id", id).Msg("banned request attempt")
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
