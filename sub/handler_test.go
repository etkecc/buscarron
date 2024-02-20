package sub

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mocks"
	"gitlab.com/etke.cc/buscarron/sub/ext/common"
)

type HandlerSuite struct {
	suite.Suite
	log    *zerolog.Logger
	v      *mocks.Validator
	vs     map[string]common.Validator
	sender *mocks.Sender
}

var ctxMatcher = mock.MatchedBy(func(ctx context.Context) bool { return true })

func (s *HandlerSuite) SetupSuite() {
	log := zerolog.Nop()
	s.log = &log
	s.v = &mocks.Validator{}
	s.vs = map[string]common.Validator{"test": s.v}
	s.sender = &mocks.Sender{}
}

func (s *HandlerSuite) SetupTest() {
	s.T().Helper()
}

func (s *HandlerSuite) TestNew() {
	handler := NewHandler(nil, s.vs, nil, s.sender, s.log)

	s.IsType(&Handler{}, handler)
}

func (s *HandlerSuite) TestGet() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{
		"test": {
			Redirect: "https://example.com",
		},
	}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)

	result, err := handler.GET(context.TODO(), "test", nil)

	s.Require().NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestGet_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)

	result, err := handler.GET(context.TODO(), "test", nil)

	s.Require().Error(err)
	s.Equal(ErrNotFound, err)
	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoForm() {
	forms := map[string]*config.Form{}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)

	result, err := handler.POST(context.TODO(), "test", nil)

	s.Require().Error(err)
	s.Equal(ErrNotFound, err)
	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_NoData() {
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com", RejectRedirect: "https://example.com"}}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)
	request, rerr := http.NewRequest("POST", "", nil)

	result, err := handler.POST(context.TODO(), "test", request)

	s.Require().NoError(rerr)
	s.Equal(ErrNotFound, err)
	s.Equal("", result)
}

func (s *HandlerSuite) TestPOST_SpamEmail() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com", RejectRedirect: "https://example.com"}}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)
	s.v.On("Email", "no").Return(false).Once()
	data := url.Values{}
	data.Add("email", "no")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST(context.TODO(), "test", request)

	s.NoError(rerr)
	s.Require().Error(err)
	s.Equal(ErrSpam, err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_SpamDomain() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com'\" /></head><body>Redirecting to <a href='https://example.com'>https://example.com</a>..."
	forms := map[string]*config.Form{"test": {Redirect: "https://example.com", RejectRedirect: "https://example.com"}}
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)
	s.v.On("Email", "").Return(true).Once()
	s.v.On("Domain", "no").Return(false).Once()
	data := url.Values{}
	data.Add("domain", "no")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST(context.TODO(), "test", request)

	s.NoError(rerr)
	s.Require().Error(err)
	s.Equal(ErrSpam, err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com/en'\" /></head><body>Redirecting to <a href='https://example.com/en'>https://example.com/en</a>..."
	// duplicated message to test extensions
	expectedMessage := "**New test** by email@dkimvalidator.com\n\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n**New test** by email@dkimvalidator.com\n\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n"
	expectedAttrs := map[string]interface{}{
		"email": "email@dkimvalidator.com",
	}
	roomID := id.RoomID("!test:example.com")
	forms := map[string]*config.Form{
		"test": {
			Name:           "test",
			Redirect:       "https://example.com/{{ .lang }}",
			RejectRedirect: "https://example.com/{{ .lang }}",
			RoomID:         roomID,
			Extensions:     []string{"", "root", "invalid"},
		},
	}
	s.v.On("Email", "email@dkimvalidator.com").Return(true).Once()
	s.v.On("Domain", "").Return(true).Once()
	s.sender.On("Send", ctxMatcher, roomID, expectedMessage, expectedAttrs).Return(id.EventID("!test:example.com")).Once()
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)
	data := url.Values{}
	data.Add("email", "email@dkimvalidator.com")
	data.Add("field", "value")
	data.Add("lang", "en")
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	result, err := handler.POST(context.TODO(), "test", request)

	s.NoError(rerr)
	s.Require().NoError(err)
	s.Equal(expected, result)
}

func (s *HandlerSuite) TestPOST_JSON() {
	expected := "<html><head><title>Redirecting...</title><meta http-equiv=\"Refresh\" content=\"0; url='https://example.com/en'\" /></head><body>Redirecting to <a href='https://example.com/en'>https://example.com/en</a>..."
	// duplicated message to test extensions
	expectedMessage := "**New test** by email@dkimvalidator.com\n\n* bool: true\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n* object: map[property:1 sub:map[sub:]]\n**New test** by email@dkimvalidator.com\n\n* bool: true\n* email: email@dkimvalidator.com\n* field: value\n* lang: en\n* object: map[property:1 sub:map[sub:]]\n"
	expectedAttrs := map[string]interface{}{
		"email": "email@dkimvalidator.com",
	}
	roomID := id.RoomID("!test:example.com")
	forms := map[string]*config.Form{
		"test": {
			Name:           "test",
			Redirect:       "https://example.com/{{ .lang }}",
			RejectRedirect: "https://example.com/{{ .lang }}",
			RoomID:         roomID,
			Extensions:     []string{"", "root", "invalid"},
		},
	}
	s.v.On("Email", "email@dkimvalidator.com").Return(true).Once()
	s.v.On("Domain", "").Return(true).Once()
	s.sender.On("Send", ctxMatcher, roomID, expectedMessage, expectedAttrs).Return(id.EventID("!test:example.com")).Once()
	handler := NewHandler(forms, s.vs, nil, s.sender, s.log)
	data := `{
	"email": "email@dkimvalidator.com",
	"field": "value",
	"bool": true,
	"object": {"property": 1, "sub": {"sub": null}},
	"lang": "en"
}`
	request, rerr := http.NewRequest("POST", "", strings.NewReader(data))
	request.Header.Add("Content-Type", "application/json")

	result, err := handler.POST(context.TODO(), "test", request)

	s.NoError(rerr)
	s.Require().NoError(err)
	s.Equal(expected, result)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
