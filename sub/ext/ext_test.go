package ext

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/etke.cc/buscarron/validator"
)

type ExtSuite struct {
	suite.Suite
}

func (s *ExtSuite) SetupTest() {
	s.T().Helper()
}

func (s *ExtSuite) TestNew() {
	v := validator.New([]string{}, []string{}, []string{}, "TRACE")
	exts := New(v, nil)

	s.IsType(&root{}, exts["root"])
}

func TestExtSuite(t *testing.T) {
	suite.Run(t, new(ExtSuite))
}
