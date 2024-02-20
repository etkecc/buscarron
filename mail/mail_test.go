package mail

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type mailSuite struct {
	suite.Suite
}

func (s *mailSuite) TestNew() {
	pm := New("test", "test@example.com", "test@example.com")

	s.IsType(&Client{}, pm)
}

func (s *mailSuite) TestNew_Empty() {
	null := New("", "", "")

	s.Nil(null)
}

func TestMail(t *testing.T) {
	suite.Run(t, new(mailSuite))
}
