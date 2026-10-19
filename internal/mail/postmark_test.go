package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/mattevans/postmark-go"
	"github.com/stretchr/testify/suite"
)

type postmarkSuite struct {
	suite.Suite
}

func (s *postmarkSuite) TestSend() {
	var mu sync.Mutex
	var path, token string
	var raw []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mu.Lock()
		path, token, raw = r.URL.Path, r.Header.Get("X-Postmark-Server-Token"), body
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"To":"to@example.com","MessageID":"msg-1"}`)
	}))
	defer srv.Close()

	err := NewPostmark("token", postmark.WithBackendURL(srv.URL)).Send(s.T().Context(), s.req())

	s.NoError(err)
	mu.Lock()
	defer mu.Unlock()
	body := map[string]any{}
	s.NoError(json.Unmarshal(raw, &body))
	s.Equal("/email", path)
	s.Equal("token", token)
	s.Equal("from@example.com", body["From"])
	s.Equal("to@example.com", body["To"])
	s.Equal("subject", body["Subject"])
	s.Equal("text body", body["TextBody"])
}

func (s *postmarkSuite) TestSend_APIError() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"ErrorCode":10,"Message":"Unprocessable"}`)
	}))
	defer srv.Close()

	err := NewPostmark("token", postmark.WithBackendURL(srv.URL)).Send(s.T().Context(), s.req())

	s.ErrorContains(err, "postmark:")
	s.ErrorContains(err, "Unprocessable")
	s.ErrorContains(err, "status: 422")
}

func (s *postmarkSuite) TestSend_TransportError() {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()

	err := NewPostmark("token", postmark.WithBackendURL(url)).Send(s.T().Context(), s.req())

	s.ErrorContains(err, "postmark:")
}

func (s *postmarkSuite) TestSend_Redirect() {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		http.Redirect(w, r, "/other", http.StatusFound)
	}))
	defer srv.Close()

	err := NewPostmark("token", postmark.WithBackendURL(srv.URL)).Send(s.T().Context(), s.req())

	s.Error(err)
	s.Equal(int32(1), hits.Load())
}

func (s *postmarkSuite) req() *postmark.Email {
	return &postmark.Email{From: "from@example.com", To: "to@example.com", Subject: "subject", TextBody: "text body"}
}

func TestPostmark(t *testing.T) {
	suite.Run(t, new(postmarkSuite))
}
