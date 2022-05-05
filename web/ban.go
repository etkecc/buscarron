package web

import (
	"net/http"
	"time"

	"git.sr.ht/~xn/cache"
	"gitlab.com/etke.cc/buscarron/logger"
)

type banhandler struct {
	store cache.Cache
	log   *logger.Logger
}

var donotban = []string{"/favicon.ico", "/robots.txt"}

// NewBanHanlder creates banhandler
func NewBanHanlder(duration int, size int, loglevel string) *banhandler {
	store := cache.NewTLRU(size, time.Duration(duration)*time.Hour, false)
	log := logger.New("ban.", loglevel)

	return &banhandler{store, log}
}

// Ban by request data
func (b *banhandler) Ban(r *http.Request) {
	for _, exclude := range donotban {
		if r.URL.Path == exclude {
			return
		}
	}
	id := r.Context().Value(ctxID)
	b.store.Set(id, true)

	b.log.Info("%s %s %v has been banned", r.Method, r.URL.Path, id)
}

// Handle to check for ban
func (b *banhandler) Handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxID)
		if b.store.Get(id) != nil {
			b.log.Debug("%s %s %v (banned) request attempt", r.Method, r.URL.String(), id)
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
