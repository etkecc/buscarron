package etkecc

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/etkecc/buscarron/internal/config"
	"github.com/etkecc/buscarron/mocks"
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

func (s *EtkeccSuite) SetupSuite() {
	s.T().Helper()
	s.v = &mocks.Validator{}
	s.ext = New(nil)
	s.ext.test = true
	s.ext.now = func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }
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
		"minimal/questions": {
			"domain":   "https://matrix.ExAmPlE.com ",
			"username": " tEsT ",
			"email":    "tEsT@TEST.cOm",
			"issue_id": "123",
		},
		"minimal/no-questions": {
			"domain":   "https://matrix.ExAmPlE.com ",
			"username": " tEsT ",
			"email":    "tEsT@TEST.cOm",
			"issue_id": "123",
			// on-premises
			"ssh-host":     "1.2.3.4",
			"ssh-port":     "222",
			"ssh-user":     "user",
			"ssh-password": "password",
			// hosting
			"ssh-client-ips": "1.2.3.4, 5.6.7.8",
			"ssh-client-key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEt3k0bEgBjfZRqU3MvWla8sgUUsm5mJRYu2CWYcYDCz user@host\nssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEt3k0bEgBjfZRqU3MvWla8sgUUsm5mJRYu2CWYcYDCz user2@host",
		},

		"full/questions": {
			"domain":              "https://matrix.ExAmPlE.com ",
			"username":            " tEsT ",
			"email":               "tEsT@TEST.cOm",
			"issue_id":            "123",
			"service-email":       "on",
			"service-support":     "dedicated",
			"baibot":              "on",
			"borg":                "on",
			"bridges-encryption":  "on",
			"buscarron":           "on",
			"cinny":               "on",
			"discord":             "on",
			"element-web":         "on",
			"etherpad":            "on",
			"facebook":            "on",
			"firezone":            "on",
			"funkwhale":           "on",
			"gmessages":           "on",
			"googlechat":          "on",
			"gotosocial":          "on",
			"honoroit":            "on",
			"hydrogen":            "on",
			"instagram":           "on",
			"irc":                 "on",
			"jitsi":               "on",
			"languagetool":        "on",
			"linkding":            "on",
			"linkedin":            "on",
			"maubot":              "on",
			"miniflux":            "on",
			"nginx-proxy-website": "on",
			"ntfy":                "on",
			"peertube":            "on",
			"postmoogle":          "on",
			"radicale":            "on",
			"reminder-bot":        "on",
			"schildichat":         "on",
			"signal":              "on",
			"skype":               "on",
			"slack":               "on",
			"sliding-sync":        "on",
			"smtp-relay":          "on",
			"stats":               "on",
			"sygnal":              "on",
			"synapse-admin":       "on",
			"synapse-s3-storage":  "on",
			"synapse-sso":         "on",
			"synapse-workers":     "on",
			"telegram":            "on",
			"twitter":             "on",
			"uptime-kuma":         "on",
			"vaultwarden":         "on",
			"webhooks":            "on",
			"wechat":              "on",
			"whatsapp":            "on",
		},
		"full/no-questions": {
			"domain":              "https://matrix.ExAmPlE.com ",
			"username":            " tEsT.admin ",
			"email":               "tEsT@TEST.cOm",
			"issue_id":            "123",
			"service-email":       "on",
			"service-support":     "dedicated",
			"baibot":              "on",
			"borg":                "on",
			"bridges-encryption":  "on",
			"buscarron":           "on",
			"cinny":               "on",
			"discord":             "on",
			"element-web":         "on",
			"etherpad":            "on",
			"facebook":            "on",
			"firezone":            "on",
			"funkwhale":           "on",
			"gmessages":           "on",
			"googlechat":          "on",
			"gotosocial":          "on",
			"honoroit":            "on",
			"hydrogen":            "on",
			"instagram":           "on",
			"irc":                 "on",
			"jitsi":               "on",
			"languagetool":        "on",
			"languagetool-ngrams": "on",
			"linkding":            "on",
			"linkedin":            "on",
			"maubot":              "on",
			"miniflux":            "on",
			"nginx-proxy-website": "on",
			"ntfy":                "on",
			"peertube":            "on",
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
			"sygnal":              "on",
			"synapse-admin":       "on",
			"synapse-s3":          "on",
			"synapse-workers":     "on",
			"telegram":            "on",
			"twitter":             "on",
			"uptime-kuma":         "on",
			"vaultwarden":         "on",
			"webhooks":            "on",
			"wechat":              "on",
			"whatsapp":            "on",
			// components
			"borg-repository":             "borg-repo",
			"funkwhale-s3-access-key":     "some-key",
			"funkwhale-s3-bucket":         "mybucket",
			"funkwhale-s3-endpoint":       "https://s3.example.com",
			"funkwhale-s3-region":         "us-east-1",
			"funkwhale-s3-secret-key":     "some-secret",
			"gotosocial-s3-access-key":    "some-key",
			"gotosocial-s3-bucket":        "mybucket",
			"gotosocial-s3-endpoint":      "https://s3.example.com",
			"gotosocial-s3-secret-key":    "some-secret",
			"nginx-proxy-website-command": "hugo",
			"nginx-proxy-website-dist":    "public",
			"nginx-proxy-website-repo":    "github.com/example/exmaple.com",
			"peertube-s3-access-key":      "some-key",
			"peertube-s3-bucket":          "mybucket",
			"peertube-s3-endpoint":        "https://s3.example.com",
			"peertube-s3-region":          "us-east-1",
			"peertube-s3-secret-key":      "some-secret",
			"reminder-bot-tz":             "America/New_York",
			"smtp-relay-email":            "user@example.com",
			"smtp-relay-host":             "smtp-relay.com",
			"smtp-relay-login":            "login",
			"smtp-relay-password":         "password",
			"smtp-relay-port":             "587",
			"sso-client-id":               "some-id",
			"sso-client-secret":           "some-secret",
			"sso-idp-brand":               "gitea",
			"sso-idp-id":                  "gitea",
			"sso-idp-name":                "gitea",
			"sso-issuer":                  "https://gitea.example.com",
			"sygnal-app-id":               "id.app",
			"sygnal-gcm-apikey":           "apikey",
			"synapse-s3-access-key":       "some-key",
			"synapse-s3-bucket":           "mybucket",
			"synapse-s3-endpoint":         "https://s3.example.com",
			"synapse-s3-region":           "us-east-1",
			"synapse-s3-secret-key":       "some-secret",
			"telegram-api-hash":           "some-hash",
			"telegram-api-id":             "123",
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
			name: "minimal/on-premises/matrix-user",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["minimal/no-questions"], map[string]string{"ssh-user": "matrix"}),
		},
		{
			name: "minimal/on-premises/domain",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/questions"],
		},
		{
			name: "minimal/on-premises/domain-a",
			before: func() {
				s.v.On("A", "example.com").Return(true).Once()
				s.v.On("NS", "example.com", "cloudflare.com").Return(true).Once()
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
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.merge(s.subs["minimal/questions"], map[string]string{"turnkey": "cpx11"}),
		},
		{
			name: "minimal/hosting/domain-a",
			before: func() {
				s.v.On("A", "example.com").Return(true).Once()
				s.v.On("NS", "example.com", "cloudflare.com").Return(true).Once()
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
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["minimal/no-questions"],
		},
		{
			name: "minimal/hosting/no-questions",
			before: func() {
				s.v.On("A", "example.com").Return(false).Once()
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
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
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
				s.v.On("CNAME", "example.com").Return(false).Once()
				s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()
			},
			submission: s.subs["full/questions"],
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
				s.v.On("NS", "example.com", "cloudflare.com").Return(false).Once()
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
func (s *EtkeccSuite) expected(name string) (questions, followup, onboarding, vars string) {
	s.T().Helper()

	return s.read(name, "questions.md"),
		s.read(name, "followup.md"),
		s.read(name, "onboarding.md"),
		s.read(name, "vars.yml")
}

// merge submissions
func (s *EtkeccSuite) merge(base, custom map[string]string) map[string]string {
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

			var actualQ, actualF, actualO, actualV string
			_, actualQ, files := s.ext.Execute(context.TODO(), s.v, &config.Form{Name: test.name}, test.submission)
			for _, file := range files {
				switch file.FileName {
				case "followup.md":
					actualF = s.rts(file.Content)
				case "onboarding.md":
					actualO = s.rts(file.Content)
				case "vars.yml":
					actualV = s.rts(file.Content)
				case "sshkey.pub", "sshkey.priv":
				default:
					s.Fail("unexpected file", file.FileName)
				}
			}
			s.saveMocks(test.name, actualQ, actualF, actualO, actualV)

			if s.save {
				s.T().Log("Save mode enabled, skipping assertions")
				return
			}
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
