package mail

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type postmarkSuite struct {
	suite.Suite
}

func (s *postmarkSuite) TestPostmarkSend() {
	s.Run("success", func() {
		var path string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.Path
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"To":"to@example.com","MessageID":"msg-1"}`))
		}))
		s.setBackend(srv.URL, srv)

		err := NewPostmark("token").Send(s.T().Context(), &postmark.Email{
			To:       "to@example.com",
			Subject:  "s",
			TextBody: "b",
		})
		s.NoError(err)
		s.Equal("/email", path)
	})

	s.Run("error", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"ErrorCode":100,"Message":"bad token"}`))
		}))
		s.setBackend(srv.URL, srv)

		err := NewPostmark("token").Send(s.T().Context(), &postmark.Email{To: "to@example.com"})
		s.Error(err)
		s.Contains(err.Error(), "bad token")
	})
}

func (s *postmarkSuite) setBackend(url string, srv *httptest.Server) {
	old := backendURL
	backendURL = url
	s.T().Cleanup(func() {
		backendURL = old
		srv.Close()
	})
}

func TestPostmark(t *testing.T) {
	suite.Run(t, new(postmarkSuite))
}
