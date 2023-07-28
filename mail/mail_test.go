package mail

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type mailSuite struct {
	suite.Suite
	log zerolog.Logger
}

func (s *mailSuite) SetupSuite() {
	s.log = zerolog.Nop()
}

func (s *mailSuite) SetupTest() {
	s.T().Helper()
}

func (s *mailSuite) TestNew() {
	pm := New("test", "test@example.com", "test@example.com", &s.log)

	s.IsType(&Client{}, pm)
}

func (s *mailSuite) TestNew_Empty() {
	null := New("", "", "", &s.log)

	s.Nil(null)
}

func TestMail(t *testing.T) {
	suite.Run(t, new(mailSuite))
}
