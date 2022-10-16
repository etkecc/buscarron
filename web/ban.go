package web

import (
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~xn/cache/v2"
	"gitlab.com/etke.cc/go/logger"
)

type banhandler struct {
	store cache.Cache[bool]
	log   *logger.Logger
}

var donotban = []string{"/favicon.ico", "/robots.txt"}

// NewBanHanlder creates banhandler
func NewBanHanlder(duration int, size int, banlist []string, loglevel string) *banhandler {
	log := logger.New("ban.", loglevel)
	store := cache.NewTLRU[bool](size, time.Duration(duration)*time.Hour, false)
	for _, id := range banlist {
		id = strings.TrimSpace(id)
		log.Info("preemptive banning %s...", id)
		store.Set(id, true)
	}

	return &banhandler{store, log}
}

// Ban by request data
func (b *banhandler) Ban(r *http.Request) {
	for _, exclude := range donotban {
		if r.URL.Path == exclude {
			return
		}
	}
	ctxIDv := r.Context().Value(ctxID)
	id, ok := ctxIDv.(string)
	if !ok {
		b.log.Error("cannot convert ctxID to string: %v", ctxIDv)
	}
	b.store.Set(id, true)

	b.log.Info("%s %s %v has been banned", r.Method, r.URL.Path, id)
}

// Handle to check for ban
func (b *banhandler) Handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxIDv := r.Context().Value(ctxID)
		id, ok := ctxIDv.(string)
		if !ok {
			b.log.Error("cannot convert ctxID to string: %v", ctxIDv)
		}
		id = strings.TrimSpace(id)
		if b.store.Has(id) {
			b.log.Info("%s %s %v (banned) request attempt", r.Method, r.URL.String(), id)
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
