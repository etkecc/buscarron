package web

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"gitlab.com/etke.cc/buscarron/logger"
)

// DefaultFrequency of the rate limiter
const DefaultFrequency = 1 * time.Minute

type ratelimiter struct {
	sync.RWMutex

	log       *logger.Logger
	burst     int
	frequency time.Duration
	visitors  map[uint32]*rlVisitor
}

// rlVisitor is a visitor's rate limiter config
type rlVisitor struct {
	*rate.Limiter
	last time.Time
}

// NewRateLimiter with defined pattern
func NewRateLimiter(pattern string, log *logger.Logger) *ratelimiter {
	burst, frequency, err := parseFrequency(pattern)
	if err != nil {
		log.Error("cannot parse rate limiter frequency pattern '%s', error: %v", pattern, err)
	}

	rl := &ratelimiter{log: log, frequency: frequency, burst: burst}
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
	default:
		frequency = DefaultFrequency
		err = fmt.Errorf("limit must be 's' or 'm' (per second or per minute), used: %s", slice[1])
	}

	return burst, frequency, err
}

func (l *ratelimiter) start() {
	ticker := time.NewTicker(l.frequency)
	for range ticker.C {
		l.log.Debug("cleanup")
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
func (l *ratelimiter) Allow(id uint32) bool {
	l.RLock()
	if l.visitors == nil {
		l.visitors = make(map[uint32]*rlVisitor)
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

	return v.Allow()
}
