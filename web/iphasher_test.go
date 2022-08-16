package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IPHasherSuite struct {
	suite.Suite
}

func (s *IPHasherSuite) TestGetHash() {
	hasher := &iphasher{}
	tests := []struct {
		input    string
		expected string
	}{
		{input: "127.0.0.1,192.168.0.1,1.1.1.1", expected: "88015183"},
		{input: "127.0.0.1", expected: "1"},
		{input: "192.168.0.1", expected: "1"},
		{input: "100.64.0.1", expected: "1"},
		{input: "1.1.1.1,1.0.0.1", expected: "87490893"},
	}

	for _, test := range tests {
		request, err := http.NewRequest("GET", "https://example.com", nil)
		request.Header.Add("X-Forwarded-For", test.input)

		hash := hasher.GetHash(request)

		s.NoError(err)
		s.Equal(test.expected, hash)
	}
}

func TestIPHasherSuite(t *testing.T) {
	suite.Run(t, new(IPHasherSuite))
}
