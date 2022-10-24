package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"maunium.net/go/mautrix/id"
)

type ConfigSuite struct {
	suite.Suite
}

var values = map[string]string{
	"BUSCARRON_LOGLEVEL": "TRACE",
	"BUSCARRON_PORT":     "12345",

	"BUSCARRON_HOMESERVER": "https://example.com",
	"BUSCARRON_LOGIN":      "@test:example.com",
	"BUSCARRON_PASSWORD":   "password",

	"BUSCARRON_SPAM_EMAILS":     "ima@spammer.com definetelynotspam@gmail.com",
	"BUSCARRON_SPAM_HOSTS":      "spamer.com unitedspammers.org",
	"BUSCARRON_SPAM_LOCALPARTS": "",
	"BUSCARRON_SPAMLIST":        "spam@*  ",

	"BUSCARRON_BAN_DURATION": "1",
	"BUSCARRON_BAN_SIZE":     "invalid",

	"BUSCARRON_LIST": "test1 test2",

	"BUSCARRON_TEST1_REDIRECT":  "https://example.org",
	"BUSCARRON_TEST1_RATELIMIT": "1r/s",
	"BUSCARRON_TEST1_ROOM":      "!test1@example.com",

	"BUSCARRON_TEST2_REDIRECT":  "https://example.com",
	"BUSCARRON_TEST2_RATELIMIT": "1r/m",
	"BUSCARRON_TEST2_ROOM":      "!test2@example.com",
}

func (s *ConfigSuite) SetupTest() {
	s.T().Helper()
	for key, value := range values {
		os.Setenv(key, value)
	}
}

func (s *ConfigSuite) TearDownTest() {
	s.T().Helper()
	for key := range values {
		os.Unsetenv(key)
	}
}

func (s *ConfigSuite) TestNew() {
	config := New()
	form1 := config.Forms["test1"]
	form2 := config.Forms["test2"]

	s.Equal("TRACE", config.LogLevel)
	s.Equal("12345", config.Port)
	s.Equal("https://example.com", config.Homeserver)
	s.Equal("@test:example.com", config.Login)
	s.Equal("password", config.Password)
	s.ElementsMatch([]string{"ima@spammer.com", "definetelynotspam@gmail.com", "*@spamer.com", "*@unitedspammers.org", "spam@*"}, config.Spamlist)
	s.Equal(1000000, config.Ban.Size)
	s.Equal("test1", form1.Name)
	s.Equal("https://example.org", form1.Redirect)
	s.Equal(id.RoomID("!test1@example.com"), form1.RoomID)
	s.Equal("1r/s", form1.Ratelimit)
	s.Equal("test2", form2.Name)
	s.Equal(id.RoomID("!test2@example.com"), form2.RoomID)
	s.Equal("https://example.com", form2.Redirect)
	s.Equal("1r/m", form2.Ratelimit)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
