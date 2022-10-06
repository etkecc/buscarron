package ext

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mocks"
)

type RootSuite struct {
	suite.Suite
	ext *root
}

func (s *RootSuite) SetupTest() {
	s.T().Helper()
	s.ext = NewRoot(&mocks.Validator{})
}

func (s *RootSuite) TestNew() {
	rootExt := NewRoot(&mocks.Validator{})

	s.IsType(&root{}, rootExt)
}

func (s *RootSuite) TestExecute() {
	expected := "**New test**\n\n* test: ✅\n\n___\n"
	data := map[string]string{
		"test": "on",
	}

	actual, _ := s.ext.Execute(&config.Form{Name: "test"}, data)

	s.Equal(expected, actual)
}

func (s *RootSuite) TestExecute_Email() {
	expected := "**New test** by test@example.com\n\n* email: test@example.com\n* test: ✅\n\n___\n"
	data := map[string]string{
		"email": "test@example.com",
		"test":  "on",
	}

	actual, _ := s.ext.Execute(&config.Form{Name: "test"}, data)

	s.Equal(expected, actual)
}

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(RootSuite))
}
