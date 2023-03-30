package web

import (
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"gitlab.com/etke.cc/buscarron/metrics"
	"gitlab.com/etke.cc/go/logger"
)

type banhandler struct {
	store *lru.Cache[string, struct{}]
	log   *logger.Logger
}

var donotban = []string{"/favicon.ico", "/robots.txt", "/metrics"}

// NewBanHanlder creates banhandler
func NewBanHanlder(size int, banlist []string, loglevel string) *banhandler {
	log := logger.New("ban.", loglevel)
	store, err := lru.New[string, struct{}](size)
	if err != nil {
		log.Error("cannot init cache: %v", err)
	}
	for _, id := range banlist {
		id = strings.TrimSpace(id)
		log.Info("preemptive banning %s...", id)
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
		b.log.Error("cannot convert ctxID to string: %v", ctxIDv)
		return
	}
	if id == "" || id == "1" {
		b.log.Error("hashed IP is empty")
		return
	}
	b.store.Add(id, struct{}{})
	ctxNameV := r.Context().Value(ctxName)
	name, ok := ctxNameV.(string)
	if !ok {
		b.log.Error("cannot convert ctxName to string: %v", ctxNameV)
	}
	metrics.BanUser(reason, name)

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
		if b.store.Contains(id) {
			metrics.BanRequest()
			b.log.Debug("%s %s %v (banned) request attempt", r.Method, r.URL.String(), id)
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
