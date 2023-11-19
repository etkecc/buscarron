package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// DefaultFrequency of the rate limiter
const DefaultFrequency = 1 * time.Minute

type rateLimiter struct {
	middlewares map[string]echo.MiddlewareFunc
	stores      map[string]*rlStore
	log         *zerolog.Logger
}

type rlStore struct {
	sync.RWMutex

	log       *zerolog.Logger
	burst     int
	frequency time.Duration
	visitors  map[string]*rlVisitor
}

// rlVisitor is a visitor's rate limiter config
type rlVisitor struct {
	*rate.Limiter
	last time.Time
}

func NewRateLimiter(shared, all map[string]string, log *zerolog.Logger) *rateLimiter {
	rl := &rateLimiter{log: log}
	rl.initStores(shared, all)
	rl.initMiddlewares()
	return rl
}

func (rl *rateLimiter) initStores(shared, all map[string]string) {
	var share *rlStore
	for _, pattern := range shared {
		if pattern == "" {
			continue
		}
		share = newStore(pattern, rl.log)
		break
	}

	stores := map[string]*rlStore{}
	for name, pattern := range all {
		if _, ok := shared[name]; ok {
			stores[name] = share
			continue
		}

		if pattern == "" {
			continue
		}
		stores[name] = newStore(pattern, rl.log)
	}
	rl.stores = stores
}

func (rl *rateLimiter) initMiddlewares() {
	mws := map[string]echo.MiddlewareFunc{}
	for name, store := range rl.stores {
		mws[name] = middleware.RateLimiter(store)
	}
	rl.middlewares = mws
}

func (rl *rateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, ok := rl.middlewares[c.Param("name")]; !ok {
				return next(c)
			}
			return rl.middlewares[c.Param("name")](next)(c)
		}
	}
}

func newStore(pattern string, log *zerolog.Logger) *rlStore {
	burst, frequency, err := parseFrequency(pattern)
	if err != nil {
		log.Error().Str("pattern", pattern).Err(err).Msg("cannot parse rate limiter frequency pattern")
	}

	rl := &rlStore{log: log, frequency: frequency, burst: burst}
	go rl.start()

	return rl
}

func parseFrequency(pattern string) (int, time.Duration, error) {
	slice := strings.Split(pattern, "r/")
	burst, err := strconv.Atoi(slice[0])
	if err != nil {
		return 1, DefaultFrequency, err
	}
	if burst < 1 {
		return 1, DefaultFrequency, fmt.Errorf("burst requests must be 1 or more, used: %d", burst)
	}

	var frequency time.Duration
	switch slice[1] {
	case "s":
		frequency = time.Duration(burst) * time.Second
	case "m":
		frequency = time.Duration(burst) * time.Minute
	case "h":
		frequency = time.Duration(burst) * time.Hour
	default:
		frequency = DefaultFrequency
		err = fmt.Errorf("limit must be 's', 'm', or 'h' (per second, per minute, or per hour), used: %s", slice[1])
	}

	return burst, frequency, err
}

func (l *rlStore) start() {
	ticker := time.NewTicker(l.frequency)
	for range ticker.C {
		l.log.Debug().Msg("cleanup")
		l.Lock()
		for id, v := range l.visitors {
			if time.Since(v.last) >= l.frequency {
				delete(l.visitors, id)
			}
		}
		l.Unlock()
	}
}

// Add new visitor
func (l *rlStore) Allow(id string) (bool, error) {
	l.RLock()
	if l.visitors == nil {
		l.visitors = make(map[string]*rlVisitor)
	}

	v, exists := l.visitors[id]
	l.RUnlock()

	if !exists {
		v = &rlVisitor{
			Limiter: rate.NewLimiter(rate.Every(l.frequency), l.burst),
		}

		l.Lock()
		l.visitors[id] = v
		l.Unlock()
	}

	v.last = time.Now()

	return v.Allow(), nil
}
