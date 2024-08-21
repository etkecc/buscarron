package ext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/etkecc/buscarron/internal/config"
)

type RootSuite struct {
	suite.Suite
	ext *root
}

func (s *RootSuite) SetupTest() {
	s.T().Helper()
	s.ext = NewRoot()
}

func (s *RootSuite) TestNew() {
	rootExt := NewRoot()

	s.IsType(&root{}, rootExt)
}

func (s *RootSuite) TestExecute() {
	expected := "**New test**\n\n* test: ✅\n"
	data := map[string]string{
		"test": "on",
	}

	actual, _ := s.ext.Execute(context.TODO(), nil, &config.Form{Name: "test"}, data)

	s.Equal(expected, actual)
}

func (s *RootSuite) TestExecute_Template() {
	expected := "**New form**:\n\ntestValue: on\n\n"
	data := map[string]string{
		"test": "on",
	}

	actual, _ := s.ext.Execute(context.TODO(), nil, &config.Form{Name: "test", Text: "**New form**:\n\ntestValue: {{ .test }}"}, data)

	s.Equal(expected, actual)
}

func (s *RootSuite) TestExecute_Email() {
	expected := "**New test** by test@example.com\n\n* email: test@example.com\n* test: ✅\n"
	data := map[string]string{
		"email": "test@example.com",
		"test":  "on",
	}

	actual, _ := s.ext.Execute(context.TODO(), nil, &config.Form{Name: "test"}, data)

	s.Equal(expected, actual)
}

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(RootSuite))
}
