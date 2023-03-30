package web

import (
	"context"
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"gitlab.com/etke.cc/go/logger"

	"gitlab.com/etke.cc/buscarron/metrics"
	"gitlab.com/etke.cc/buscarron/sub"
)

// FormHandler for web server
type FormHandler interface {
	GET(string, *http.Request) (string, error)
	POST(string, string, *http.Request) (string, error)
}

// DomainValidator for the web server
type DomainValidator interface {
	A(string) bool
	DomainString(string) bool
	GetBase(string) string
}

// Server to handle forms
type Server struct {
	fh  FormHandler
	dv  DomainValidator
	sh  *sentryhttp.Handler
	bh  *banhandler
	log *logger.Logger
	iph *iphasher
	rls map[string]*ratelimiter
	frr map[string]string
	srv *http.Server
}

// New web server
func New(port string, srl, rls map[string]string, frr map[string]string, loglevel string, fh FormHandler, dv DomainValidator, bs int, bl []string) *Server {
	log := logger.New("web.", loglevel)
	sh := sentryhttp.New(sentryhttp.Options{})
	bh := NewBanHanlder(bs, bl, loglevel)
	iph := &iphasher{}
	ctxm := &ctxMiddleware{iph}
	srv := &Server{
		bh:  bh,
		dv:  dv,
		fh:  fh,
		sh:  sh,
		log: log,
		iph: &iphasher{},
		frr: frr,
		rls: initRateLimiters(srl, rls, log),
	}

	mux := http.NewServeMux()
	metrics.InitMetrics(mux)
	mux.Handle("/_health", srv.healthcheck())
	mux.Handle("/_domain", srv.domainValidator())
	mux.Handle("/", srv.forms())

	h := cors(ctxm.Handle(sh.Handle(bh.Handle(mux))))
	srv.srv = &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	return srv
}

func initRateLimiters(srl, rlCfg map[string]string, log *logger.Logger) map[string]*ratelimiter {
	var shared *ratelimiter
	for _, cfg := range srl {
		if cfg == "" {
			continue
		}
		shared = NewRateLimiter(cfg, log)
		break
	}

	rls := map[string]*ratelimiter{}
	for name, cfg := range rlCfg {
		if _, ok := srl[name]; ok {
			rls[name] = shared
			continue
		}

		if cfg == "" {
			continue
		}
		rls[name] = NewRateLimiter(cfg, log)
	}

	return rls
}

func (s *Server) healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			s.log.Error("%s %s %v", r.Method, r.URL.String(), err)
		}
	}
}

func (s *Server) domainValidator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		domain = s.dv.GetBase(domain)

		if s.dv.DomainString(domain) && !s.dv.A(domain) {
			s.log.Info("%s %s %s is valid", r.Method, r.URL.String(), domain)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		s.log.Info("%s %s %s is invalid", r.Method, r.URL.String(), domain)
		http.Error(w, "", http.StatusConflict)
	}
}

func (s *Server) forms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxID).(string)
		name := r.Context().Value(ctxName).(string)

		if r.Method == http.MethodPost {
			s.formPOST(id, name, w, r)
			return
		}

		s.formGET(name, w, r)
	}
}

func (s *Server) formReject(name, reason string, w http.ResponseWriter, r *http.Request) {
	target, ok := s.frr[name]
	s.log.Info("%s %s submission to %s rejected, reason: %s", r.Method, r.URL.Path, name, reason)
	if !ok {
		return
	}
	w.Write([]byte(`<html><head><title>Redirecting...</title><meta http-equiv="Refresh" content="0; url='` + target + `'" /></head><body>Redirecting to <a href='` + target + `'>` + target + `</a>...`)) //nolint:errcheck
}

func (s *Server) formGET(name string, w http.ResponseWriter, r *http.Request) {
	body, err := s.fh.GET(name, r)
	if err == sub.ErrNotFound || err == sub.ErrSpam {
		s.bh.Ban(r, err.Error())
		s.formReject(name, err.Error(), w, r)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		s.log.Error("%s %s %v", r.Method, r.URL.Path, err)
	}
}

func (s *Server) formPOST(id, name string, w http.ResponseWriter, r *http.Request) {
	var limited bool
	rl := s.rls[name]
	if rl != nil {
		limited = !rl.Allow(id)
	}

	if limited {
		http.Error(w, "", http.StatusTooManyRequests)
		s.formReject(name, "too many requests", w, r)
		return
	}

	body, err := s.fh.POST(id, name, r)
	if err == sub.ErrNotFound || err == sub.ErrSpam {
		s.bh.Ban(r, err.Error())
		s.formReject(name, err.Error(), w, r)
		return
	}

	if _, err := w.Write([]byte(body)); err != nil {
		s.log.Error("%s %s %v", r.Method, r.URL.Path, err)
	}
}

// Start web server
func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop web server
func (s *Server) Stop() {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		s.log.Error("cannot stop web server: %v", err)
	}
}
