package web

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/time/rate"

	"gitlab.com/etke.cc/buscarron/logger"
)

type RatelimiterSuite struct {
	suite.Suite
}

func (s *RatelimiterSuite) TestNewRateLimiter() {
	s.T().Parallel()
	log := logger.New("rl.", "TRACE")
	tests := []struct {
		input string
		burst int
		limit rate.Limit
	}{
		{input: "1r/s", burst: 1, limit: rate.Every(1 * time.Second)},
		{input: "3r/m", burst: 3, limit: rate.Every(1 * time.Minute)},
		{input: "invalidr/h", burst: 1, limit: rate.Every(DefaultFrequency)},
		{input: "0r/h", burst: 1, limit: rate.Every(DefaultFrequency)},
		{input: "5r/h", burst: 5, limit: rate.Every(DefaultFrequency)},
	}

	for _, test := range tests {
		s.Run(test.input, func() {
			rl := NewRateLimiter(test.input, log)

			s.Equal(test.burst, rl.burst)
		})
	}
}

func (s *RatelimiterSuite) TestAllow() {
	s.T().Parallel()
	log := logger.New("rl.", "TRACE")
	rl := NewRateLimiter("1r/s", log)

	first := rl.Allow(1)
	second := rl.Allow(1)
	time.Sleep(2 * time.Second)
	third := rl.Allow(1)
	fourth := rl.Allow(1)

	s.True(first)
	s.False(second)
	s.True(third)
	s.False(fourth)
}

func TestRatelimiterSuite(t *testing.T) {
	suite.Run(t, new(RatelimiterSuite))
}
