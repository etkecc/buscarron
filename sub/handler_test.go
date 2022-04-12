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
	handler := NewHandler(nil, &config.Spam{}, s.sender, "TRACE")

	s.IsType(&Handler{}, handler)
}

func (s *HandlerSuite) TestGet() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{
		"test": {
			Redirect: "https://example.com",
		},
	}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")

	result := handler.GET("test", nil)

	s.Equal(expected, result)
}

func (s *HandlerSuite) TestGet_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")

	result := handler.GET("test", nil)

	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")

	result := handler.POST("test", nil)

	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoData() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")
	request, err := http.NewRequest("POST", "", nil)

	result := handler.POST("test", request)

	s.NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_SpamEmail() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")
	data := url.Values{}
	data.Add("email", "no")
	request, err := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result := handler.POST("test", request)

	s.NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_SpamDomain() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com"}}
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")
	data := url.Values{}
	data.Add("domain", "no")
	request, err := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result := handler.POST("test", request)

	s.NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	// duplicated message to test extensions
	expectedMessage := "**New test** by test@example.com\n\n* email: test@example.com\n* field: value\n\n___\n**New test** by test@example.com\n\n* email: test@example.com\n* field: value\n\n___\n"
	roomID := id.RoomID("!test:example.com")
	forms := map[string]*config.Form{
		"test": {
			Name:       "test",
			Redirect:   "https://example.com",
			RoomID:     roomID,
			Extensions: []string{"", "root", "invalid"},
		},
	}
	s.sender.On("Send", roomID, expectedMessage).Once()
	handler := NewHandler(forms, &config.Spam{}, s.sender, "TRACE")
	data := url.Values{}
	data.Add("email", "test@example.com")
	data.Add("field", "value")
	request, err := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result := handler.POST("test", request)

	s.NoError(err)
	s.Equal(expected, result)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
