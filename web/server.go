package web

import (
	"context"
	"net/http"
	"strings"

	"gitlab.com/etke.cc/buscarron/logger"
)

// FormHandler for web server
type FormHandler interface {
	GET(string, *http.Request) string
	POST(string, *http.Request) string
}

// Server to handle forms
type Server struct {
	fh   FormHandler
	log  *logger.Logger
	iph  *iphasher
	rls  map[string]*ratelimiter
	http *http.Server
}

// New web server
func New(port string, rls map[string]string, loglevel string, fh FormHandler) *Server {
	log := logger.New("web.", loglevel)
	srv := &Server{
		fh:  fh,
		log: log,
		iph: &iphasher{},
		rls: initRateLimiters(rls, log),
		http: &http.Server{
			Addr:     ":" + port,
			ErrorLog: log.GetLog(),
		},
	}
	srv.initHealthcheck()
	srv.initForms()

	return srv
}

func initRateLimiters(rlCfg map[string]string, log *logger.Logger) map[string]*ratelimiter {
	rls := map[string]*ratelimiter{}
	for name, cfg := range rlCfg {
		if cfg == "" {
			continue
		}
		rls[name] = NewRateLimiter(cfg, log)
	}

	return rls
}

func (s *Server) initHealthcheck() {
	http.HandleFunc("/_health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			s.log.Error("%s %s %v", r.Method, r.URL.String(), err)
		}
	})
}

func (s *Server) initForms() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		name := strings.ReplaceAll(url, "/", "")
		method := r.Method
		if method != http.MethodPost {
			body := s.fh.GET(name, r)
			if _, err := w.Write([]byte(body)); err != nil {
				s.log.Error("%s %s %v", method, url, err)
			}
			return
		}

		if s.isLimited(name, r) {
			http.Error(w, "", http.StatusTooManyRequests)
			s.log.Error("%s %s too many requests", method, url)
			return
		}

		body := s.fh.POST(name, r)
		if _, err := w.Write([]byte(body)); err != nil {
			s.log.Error("%s %s %v", method, url, err)
		}
	})
}

func (s *Server) isLimited(name string, r *http.Request) bool {
	id := s.iph.GetHash(r)
	rl, ok := s.rls[name]
	if rl == nil || !ok {
		return false
	}

	return !rl.Allow(id)
}

// Start web server
func (s *Server) Start() error {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop web server
func (s *Server) Stop() {
	if err := s.http.Shutdown(context.Background()); err != nil {
		s.log.Error("cannot stop web server: %v", err)
	}
}
