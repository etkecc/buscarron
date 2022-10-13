package ext

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExtSuite struct {
	suite.Suite
}

func (s *ExtSuite) SetupTest() {
	s.T().Helper()
}

func (s *ExtSuite) TestNew() {
	exts := New(nil)

	s.IsType(&root{}, exts["root"])
}

func TestExtSuite(t *testing.T) {
	suite.Run(t, new(ExtSuite))
}
