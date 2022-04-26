package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"gitlab.com/etke.cc/buscarron/logger"
)

// FormHandler for web server
type FormHandler interface {
	GET(string, *http.Request) string
	POST(string, *http.Request) string
}

// DomainValidator for the web server
type DomainValidator interface {
	A(string) bool
	DomainString(string) bool
}

// Server to handle forms
type Server struct {
	fh   FormHandler
	dv   DomainValidator
	sh   *sentryhttp.Handler
	log  *logger.Logger
	iph  *iphasher
	rls  map[string]*ratelimiter
	http *http.Server
}

// New web server
func New(port string, rls map[string]string, loglevel string, fh FormHandler, dv DomainValidator) *Server {
	log := logger.New("web.", loglevel)
	sh := sentryhttp.New(sentryhttp.Options{})
	srv := &Server{
		fh:  fh,
		dv:  dv,
		sh:  sh,
		log: log,
		iph: &iphasher{},
		rls: initRateLimiters(rls, log),
		http: &http.Server{
			Addr:     ":" + port,
			ErrorLog: log.GetLog(),
		},
	}
	srv.initHealthcheck()
	srv.initDomainValidator()
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
	http.HandleFunc("/_health", s.sh.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			s.log.Error("%s %s %v", r.Method, r.URL.String(), err)
		}
	}))
}

func (s *Server) initDomainValidator() {
	http.HandleFunc("/_domain", s.sh.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		if s.dv.DomainString(domain) && !s.dv.A(domain) {
			s.log.Info("%s %s %s is valid", r.Method, r.URL.String(), domain)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		s.log.Info("%s %s %s is invalid", r.Method, r.URL.String(), domain)
		http.Error(w, "", http.StatusForbidden)
	}))
}

func (s *Server) initForms() {
	http.HandleFunc("/", s.sh.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		name := strings.ReplaceAll(url, "/", "")
		method := r.Method
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("form", "name")
			})
		}
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
	}))
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
