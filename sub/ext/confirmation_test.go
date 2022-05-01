package ext

import (
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mocks"
)

type ConfirmationSuite struct {
	suite.Suite
	ext *confirmation
	s   *mocks.EmailSender
}

func (s *ConfirmationSuite) SetupTest() {
	s.T().Helper()
	s.s = &mocks.EmailSender{}
	s.ext = NewConfirmation(s.s)
}

func (s *ConfirmationSuite) TearDownTest() {
	s.T().Helper()
	s.s.AssertExpectations(s.T())
}

func (s *ConfirmationSuite) TestNew() {
	confirmationExt := NewConfirmation(s.s)

	s.IsType(&confirmation{}, confirmationExt)
}

func (s *ConfirmationSuite) TestExecute() {
	form := &config.Form{
		Name: "test",
		Confirmation: config.Confirmation{
			Subject: "Confirmation email with testfield = {{ .testfield }}",
			Body:    "With body, too ({{ .testfield }})",
		},
	}
	data := map[string]string{
		"email":     "test@example.com",
		"testfield": "testvalue",
	}
	req := &postmark.Email{
		To:       "test@example.com",
		Tag:      "confirmation",
		Subject:  "Confirmation email with testfield = testvalue",
		TextBody: "With body, too (testvalue)",
	}
	s.s.On("Send", req).Return(nil, nil).Once()

	s.ext.Execute(form, data)
}

func (s *ConfirmationSuite) TestExecute_NoPostmark() {
	ext := NewConfirmation(nil)
	form := &config.Form{
		Name: "test",
		Confirmation: config.Confirmation{
			Subject: "Confirmation email with testfield = {{ .testfield }}",
			Body:    "With body, too ({{ .testfield }})",
		},
	}
	data := map[string]string{
		"email":     "test@example.com",
		"testfield": "testvalue",
	}

	ext.Execute(form, data)
}

func (s *ConfirmationSuite) TestExecute_NotConfigured() {
	form := &config.Form{
		Name: "test",
	}
	data := map[string]string{
		"email":     "test@example.com",
		"testfield": "testvalue",
	}

	s.ext.Execute(form, data)
}

func (s *ConfirmationSuite) TestExecute_NoEmail() {
	form := &config.Form{
		Name: "test",
		Confirmation: config.Confirmation{
			Subject: "Confirmation email with testfield = {{ .testfield }}",
			Body:    "With body, too ({{ .testfield }})",
		},
	}
	data := map[string]string{
		"testfield": "testvalue",
	}

	s.ext.Execute(form, data)
}

func (s *ConfirmationSuite) TestExecute_ErrorParsingSubject() {
	form := &config.Form{
		Name: "test",
		Confirmation: config.Confirmation{
			Subject: "Confirmation email with testfield = {{ {{ $@$@$@testfield }}",
			Body:    "With body, too ({{ .testfield }})",
		},
	}
	data := map[string]string{
		"email":     "test@example.com",
		"testfield": "testvalue",
	}

	s.ext.Execute(form, data)
}

func (s *ConfirmationSuite) TestExecute_ErrorParsingBody() {
	form := &config.Form{
		Name: "test",
		Confirmation: config.Confirmation{
			Body: "Confirmation email with testfield = {{ {{ $@$@$@testfield }}",
		},
	}
	data := map[string]string{
		"email":     "test@example.com",
		"testfield": "testvalue",
	}

	s.ext.Execute(form, data)
}

func TestConfirmationSuite(t *testing.T) {
	suite.Run(t, new(ConfirmationSuite))
}
