package sub

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mocks"
	"gitlab.com/etke.cc/buscarron/validator"
)

type HandlerSuite struct {
	suite.Suite
	sender *mocks.Sender
}

func (s *HandlerSuite) SetupTest() {
	s.T().Helper()
	s.sender = &mocks.Sender{}
}

func (s *HandlerSuite) TestNew() {
	handler := NewHandler(nil, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")

	s.IsType(&Handler{}, handler)
}

func (s *HandlerSuite) TestGet() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{
		"test": {
			Redirect: "https://example.com",
		},
	}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")

	result, err := handler.GET("test", nil)

	s.Require().NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestGet_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")

	result, err := handler.GET("test", nil)

	s.Require().Error(err)
	s.Equal(ErrNotFound, err)
	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")

	result, err := handler.POST("test", nil)

	s.Require().Error(err)
	s.Equal(ErrNotFound, err)
	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoData() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")
	request, rerr := http.NewRequest("POST", "", nil)

	result, err := handler.POST("test", request)

	s.Require().NoError(rerr)
	s.Require().NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_SpamEmail() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")
	data := url.Values{}
	data.Add("email", "no")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST("test", request)

	s.NoError(rerr)
	s.Require().Error(err)
	s.Equal(ErrSpam, err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_SpamDomain() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")
	data := url.Values{}
	data.Add("domain", "no")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST("test", request)

	s.NoError(rerr)
	s.Require().Error(err)
	s.Equal(ErrSpam, err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com/en'\" /></head><body>Redirecting to <a href='https://example.com/en'>https://example.com/en</a>..."
	// duplicated message to test extensions
	expectedMessage := "**New test** by email@dkimvalidator.com\n\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n\n___\n**New test** by email@dkimvalidator.com\n\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n\n___\n"
	roomID := id.RoomID("!test:example.com")
	forms := map[string]*config.Form{
		"test": {
			Name:       "test",
			Redirect:   "https://example.com/{{ .lang }}",
			RoomID:     roomID,
			Extensions: []string{"", "root", "invalid"},
		},
	}
	s.sender.On("Send", roomID, expectedMessage).Once()
	handler := NewHandler(forms, validator.New([]string{}, []string{}, []string{}, "TRACE"), nil, s.sender, "TRACE")
	data := url.Values{}
	data.Add("email", "email@dkimvalidator.com")
	data.Add("field", "value")
	data.Add("lang", "en")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST("test", request)

	s.NoError(rerr)
	s.Require().NoError(err)
	s.Equal(expected, result)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
