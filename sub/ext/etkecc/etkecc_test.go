package etkecc

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/etke.cc/buscarron/config"
	"gitlab.com/etke.cc/buscarron/mocks"
)

type testCase struct {
	name       string
	before     func()
	submission map[string]string
}

type EtkeccSuite struct {
	suite.Suite
	v     *mocks.Validator
	ext   *Etkecc
	save  bool
	subs  map[string]map[string]string
	cases []testCase
}

func (s *EtkeccSuite) SetupTest() {
	s.T().Helper()
	s.v = &mocks.Validator{}
	s.ext = New(nil)
	s.ext.test = true
	s.save = false

	s.setupSubs()
	s.setupCases()
}

func (s *EtkeccSuite) TearDownTest() {
	s.T().Helper()
	s.v.AssertExpectations(s.T())
}

// setupSubs inits base submissions
func (s *EtkeccSuite) setupSubs() {
	s.T().Helper()

	s.subs = map[string]map[string]string{
		// OLD forms
		"minimal/old/on-premises": {
			"domain":   "https://matrix.ExAmPlE.com ",
			"username": " tEsT ",
			"email":    "tEsT@TEST.cOm",
			"type":     "byos",
		},
		"minimal/old/hosting": {
			"domain":           "https://matrix.ExAmPlE.com ",
			"username":         " tEsT ",
			"email":            "tEsT@TEST.cOm",
			"type":             "turnkey",
			"turnkey":          "small-cx11",
			"turnkey-location": "Nuremberg",
		},

		"minimal/questions": {
			"domain":   "https://matrix.ExAmPlE.com ",
			"username": " tEsT ",
			"email":    "tEsT@TEST.cOm",
		},
		"minimal/no-questions": {
			"domain":   "https://matrix.ExAmPlE.com ",
			"username": " tEsT ",
			"email":    "tEsT@TEST.cOm",
			// on-premises
			"ssh-host":     "1.2.3.4",
			"ssh-port":     "222",
			"ssh-user":     "user",
			"ssh-password": "password",
			// hosting
			"ssh-client-ips": "1.2.3.4, 5.6.7.8",
			"ssh-client-key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEt3k0bEgBjfZRqU3MvWla8sgUUsm5mJRYu2CWYcYDCz user@host",
		},

		"full/questions": {
			"domain":              "https://matrix.ExAmPlE.com ",
			"username":            " tEsT ",
			"email":               "tEsT@TEST.cOm",
			"service-email":       "yes",
			"service-support":     "dedicated",
			"borg":                "on",
			"bridges-encryption":  "on",
			"buscarron":           "on",
			"chatgpt":             "on",
			"cinny":               "on",
			"discord":             "on",
			"element-web":         "on",
			"etherpad":            "on",
			"facebook":            "on",
			"gmessages":           "on",
			"googlechat":          "on",
			"gotosocial":          "on",
			"groupme":             "on",
			"honoroit":            "on",
			"hydrogen":            "on",
			"instagram":           "on",
			"irc":                 "on",
			"jitsi":               "on",
			"linkedin":            "on",
			"miniflux":            "on",
			"nginx-proxy-website": "on",
			"ntfy":                "on",
			"postmoogle":          "on",
			"radicale":            "on",
			"reminder-bot":        "on",
			"schildichat":         "on",
			"signal":              "on",
			"skype":               "on",
			"slack":               "on",
			"sliding-sync":        "on",
			"smtp-relay":          "on",
			"sso":                 "on",
			"stats":               "on",
			"steam":               "on",
			"sygnal":              "on",
			"synapse-admin":       "on",
			"telegram":            "on",
			"twitter":             "on",
			"uptime-kuma":         "on",
			"vaultwarden":         "on",
			"webhooks":            "on",
			"whatsapp":            "on",
		},
		"full/no-questions": {
			"domain":              "https://matrix.ExAmPlE.com ",
			"username":            " tEsT ",
			"email":               "tEsT@TEST.cOm",
			"service-email":       "yes",
			"service-support":     "dedicated",
			"borg":                "on",
			"bridges-encryption":  "on",
			"buscarron":           "on",
			"chatgpt":             "on",
			"cinny":               "on",
			"discord":             "on",
			"element-web":         "on",
			"etherpad":            "on",
			"facebook":            "on",
			"gmessages":           "on",
			"googlechat":          "on",
			"gotosocial":          "on",
			"groupme":             "on",
			"honoroit":            "on",
			"hydrogen":            "on",
			"instagram":           "on",
			"irc":                 "on",
			"jitsi":               "on",
			"linkedin":            "on",
			"miniflux":            "on",
			"nginx-proxy-website": "on",
			"ntfy":                "on",
			"postmoogle":          "on",
			"radicale":            "on",
			"reminder-bot":        "on",
			"schildichat":         "on",
			"signal":              "on",
			"skype":               "on",
			"slack":               "on",
			"sliding-sync":        "on",
			"smtp-relay":          "on",
			"sso":                 "on",
			"stats":               "on",
			"steam":               "on",
			"sygnal":              "on",
			"synapse-admin":       "on",
			"telegram":            "on",
			"twitter":             "on",
			"uptime-kuma":         "on",
			"vaultwarden":         "on",
			"webhooks":            "on",
			"whatsapp":            "on",
			// components
			"reminder-bot-tz":             "America/New_York",
			"telegram-api-id":             "123",
			"telegram-api-hash":           "some-hash",
			"smtp-relay-host":             "smtp-relay.com",
			"smtp-relay-port":             "587",
			"smtp-relay-login":            "login",
			"smtp-relay-password":         "password",
			"smtp-relay-email":            "user@example.com",
			"nginx-proxy-website-command": "hugo",
			"nginx-proxy-website-dist":    "public",
			"nginx-proxy-website-repo":    "github.com/example/exmaple.com",
			"sso-client-id":               "some-id",
			"sso-client-secret":           "some-secret",
			"sso-issuer":                  "example",
			"sso-idp-brand":               "example",
			"sso-idp-id":                  "example",
			"sso-idp-name":                "example",
			"sygnal-app-id":               "id.app",
			"sygnal-gcm-apikey":           "apikey",
			"borg-repository":             "borg-repo",
			// on-premises
			"ssh-host":     "1.2.3.4",
			"ssh-port":     "222",
			"ssh-user":     "user",
			"ssh-password": "password",
			// hosting
			"ssh-client-ips": "1.2.3.4, 5.6.7.8",
			"ssh-client-key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEt3k0bEgBjfZRqU3MvWla8sgUUsm5mJRYu2CWYcYDCz user@host",
		},
	}
}

// setupCases inits test cases
func (s *EtkeccSuite) setupCases() {
	s.T().Helper()
	s.cases = []testCase{
		{
			name: "minimal/on-premises/old",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/old/on-premises"],
		},
		{
			name: "minimal/hosting/old",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/old/hosting"],
		},

		{
			name: "minimal/on-premises/domain",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/questions"],
		},
		{
			name: "minimal/on-premises/domain-a",
			before: func() {
				s.v.On("A", "example.com").Return(true).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/questions"],
		},
		{
			name: "minimal/on-premises/subdomain",
			before: func() {
				s.v.On("A", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("CNAME", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("GetBase", "https://higenjitsuteki.onmatrix.chat").Return("higenjitsuteki.onmatrix.chat").Once()
			},
			submission: s.merge(s.subs["minimal/questions"], map[string]string{"domain": "https://higenjitsuteki.onmatrix.chat"}),
		},
		{
			name: "minimal/hosting/domain",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["minimal/questions"], map[string]string{"turnkey": "cpx11"}),
		},
		{
			name: "minimal/hosting/subdomain",
			before: func() {
				s.v.On("A", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("CNAME", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("GetBase", "https://higenjitsuteki.onmatrix.chat").Return("higenjitsuteki.onmatrix.chat").Once()
			},
			submission: s.merge(s.subs["minimal/questions"], map[string]string{
				"domain":  "https://higenjitsuteki.onmatrix.chat",
				"turnkey": "cpx11",
			}),
		},
		{
			name: "minimal/on-premises/no-questions",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/no-questions"],
		},
		{
			name: "minimal/hosting/no-questions",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["minimal/no-questions"], map[string]string{
				"turnkey": "cpx11",
			}),
		},
		{
			name: "full/on-premises/questions",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["full/questions"], map[string]string{"service-email": "no"}),
		},
		{
			name: "full/hosting/questions",
			before: func() {
				s.v.On("A", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("CNAME", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("GetBase", "https://higenjitsuteki.onmatrix.chat").Return("higenjitsuteki.onmatrix.chat").Once()
			},
			submission: s.merge(s.subs["full/questions"], map[string]string{"turnkey": "cpx11", "domain": "https://higenjitsuteki.onmatrix.chat"}),
		},
		{
			name: "full/on-premises/no-questions",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["full/no-questions"], map[string]string{"service-email": "no"}),
		},
		{
			name: "full/hosting/no-questions",
			before: func() {
				s.v.On("A", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("CNAME", "higenjitsuteki.onmatrix.chat").Return(false).Once()
				s.v.On("GetBase", "https://higenjitsuteki.onmatrix.chat").Return("higenjitsuteki.onmatrix.chat").Once()
			},
			submission: s.merge(s.subs["full/no-questions"], map[string]string{"turnkey": "cpx11", "domain": "https://higenjitsuteki.onmatrix.chat"}),
		},
	}
}

// expected returns expected questions, followup, onboarding, vars
func (s *EtkeccSuite) expected(name string) (string, string, string, string) {
	s.T().Helper()

	return s.read(name, "questions.md"),
		s.read(name, "followup.md"),
		s.read(name, "onboarding.md"),
		s.read(name, "vars.yml")
}

// merge submissions
func (s *EtkeccSuite) merge(base map[string]string, custom map[string]string) map[string]string {
	s.T().Helper()
	merged := map[string]string{}
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}

	return merged
}

// saveMocks stores the actual generated data if s.save = true
func (s *EtkeccSuite) saveMocks(name, questions, followup, onboarding, vars string) {
	s.T().Helper()
	if !s.save {
		return
	}

	s.write(name, "questions.md", questions)
	s.write(name, "followup.md", followup)
	s.write(name, "onboarding.md", onboarding)
	s.write(name, "vars.yml", vars)
}

// read file from the fs
func (s *EtkeccSuite) read(name, file string) string {
	s.T().Helper()

	text, _ := os.ReadFile("testdata/" + name + "/" + file)

	return string(text)
}

// write file to the fs
func (s *EtkeccSuite) write(name, file, content string) {
	s.T().Helper()

	err := os.MkdirAll("testdata/"+name, 0o700)
	s.NoError(err)

	err = os.WriteFile("testdata/"+name+"/"+file, []byte(content), 0o666)
	s.NoError(err)
}

// rts is io.Reader to string
func (s *EtkeccSuite) rts(r io.Reader) string {
	s.T().Helper()

	var buf strings.Builder
	_, err := io.Copy(&buf, r)
	s.NoError(err)

	return buf.String()
}

func (s *EtkeccSuite) TestNew() {
	etkeccExt := New(nil)

	s.IsType(&Etkecc{}, etkeccExt)
}

func (s *EtkeccSuite) TestPassword() {
	expected := "value"
	o := &order{pass: map[string]string{}}
	o.pass["test"] = "value"

	actual := o.password("test")

	s.Equal(expected, actual)
}

func (s *EtkeccSuite) TestExecute() {
	for _, test := range s.cases {
		s.Run(test.name, func() {
			expectedQ, expectedF, expectedO, expectedV := s.expected(test.name)
			test.before()

			actualQ, files := s.ext.Execute(s.v, &config.Form{Name: test.name}, test.submission)
			actualV := s.rts(files[0].Content)
			actualO := s.rts(files[1].Content)
			actualF := s.rts(files[2].Content)
			s.saveMocks(test.name, actualQ, actualF, actualO, actualV)

			s.Equal(expectedQ, actualQ)
			s.Equal(expectedF, actualF)
			s.Equal(expectedO, actualO)
			s.Equal(expectedV, actualV)
		})
	}
}

func TestEtkeccSuite(t *testing.T) {
	suite.Run(t, new(EtkeccSuite))
}
