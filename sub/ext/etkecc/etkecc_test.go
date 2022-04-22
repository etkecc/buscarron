package etkecc

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/etke.cc/buscarron/mocks"
)

type EtkeccSuite struct {
	suite.Suite
	v           *mocks.NetworkValidator
	ext         *Etkecc
	save        bool
	byos        map[string]string
	turnkey     map[string]string
	byosFull    map[string]string
	turnkeyFull map[string]string
}

func (s *EtkeccSuite) SetupTest() {
	s.T().Helper()
	s.v = &mocks.NetworkValidator{}
	s.ext = New(s.v)
	s.ext.test = true
	s.save = false

	s.byos = map[string]string{
		"homeserver": "synapse",
		"domain":     "https://matrix.ExAmPlE.com ",
		"username":   " tEsT ",
		"email":      "tEsT@TEST.cOm",
		"type":       "byos",
	}
	s.turnkey = map[string]string{
		"homeserver": "synapse",
		"domain":     "https://matrix.ExAmPlE.com ",
		"username":   " tEsT ",
		"email":      "tEsT@TEST.cOm",
		"type":       "turnkey",
	}
	s.byosFull = map[string]string{
		"homeserver":              "synapse",
		"domain":                  "https://matrix.ExAmPlE.com ",
		"username":                " tEsT ",
		"email":                   "tEsT@TEST.cOm",
		"name":                    "Test",
		"notes":                   "ALL, I.WANT.THEM.ALLLLLLL. Yes, on shitty 1vCPU 1GB RAM VPS.",
		"service-setup":           "on",
		"service-maintenance":     "on",
		"service-email":           "on",
		"type":                    "byos",
		"borg":                    "on",
		"buscarron":               "on",
		"cinny":                   "on",
		"dimension":               "on",
		"discord":                 "on",
		"dnsmasq":                 "on",
		"element-web":             "on",
		"email2matrix":            "on",
		"etherpad":                "on",
		"facebook":                "on",
		"go-neb":                  "on",
		"googlechat":              "on",
		"groupme":                 "on",
		"hangouts":                "on",
		"heisenbridge":            "on",
		"honoroit":                "on",
		"hydrogen":                "on",
		"instagram":               "on",
		"irc":                     "on",
		"jitsi":                   "on",
		"kuma":                    "on",
		"languagetool":            "on",
		"linkedin":                "on",
		"ma1sd":                   "on",
		"matrix-corporal":         "on",
		"matrix-registration":     "on",
		"miniflux":                "on",
		"miounne":                 "on",
		"mjolnir":                 "on",
		"nginx-proxy-website":     "on",
		"radicale":                "on",
		"reminder-bot":            "on",
		"signal":                  "on",
		"skype":                   "on",
		"slack":                   "on",
		"smtp-relay":              "on",
		"sso":                     "on",
		"stats":                   "on",
		"steam":                   "on",
		"sygnal":                  "on",
		"synapse-simple-antispam": "on",
		"synapse-workers":         "on",
		"telegram":                "on",
		"twitter":                 "on",
		"uptime-kuma":             "on",
		"webhooks":                "on",
		"whatsapp":                "on",
		"wireguard":               "on",
	}
	s.turnkeyFull = map[string]string{
		"homeserver":              "synapse",
		"domain":                  "https://matrix.ExAmPlE.com ",
		"username":                " tEsT ",
		"email":                   "tEsT@TEST.cOm",
		"name":                    "Test",
		"notes":                   "ALL, I.WANT.THEM.ALLLLLLL. Yes, on shitty 1vCPU 1GB RAM VPS.",
		"service-setup":           "on",
		"service-maintenance":     "on",
		"service-email":           "on",
		"type":                    "turnkey",
		"borg":                    "on",
		"buscarron":               "on",
		"cinny":                   "on",
		"dimension":               "on",
		"discord":                 "on",
		"dnsmasq":                 "on",
		"element-web":             "on",
		"email2matrix":            "on",
		"etherpad":                "on",
		"facebook":                "on",
		"go-neb":                  "on",
		"googlechat":              "on",
		"groupme":                 "on",
		"hangouts":                "on",
		"heisenbridge":            "on",
		"honoroit":                "on",
		"hydrogen":                "on",
		"instagram":               "on",
		"irc":                     "on",
		"jitsi":                   "on",
		"kuma":                    "on",
		"languagetool":            "on",
		"linkedin":                "on",
		"ma1sd":                   "on",
		"matrix-corporal":         "on",
		"matrix-registration":     "on",
		"miniflux":                "on",
		"miounne":                 "on",
		"mjolnir":                 "on",
		"nginx-proxy-website":     "on",
		"radicale":                "on",
		"reminder-bot":            "on",
		"signal":                  "on",
		"skype":                   "on",
		"slack":                   "on",
		"smtp-relay":              "on",
		"sso":                     "on",
		"stats":                   "on",
		"steam":                   "on",
		"sygnal":                  "on",
		"synapse-simple-antispam": "on",
		"synapse-workers":         "on",
		"telegram":                "on",
		"twitter":                 "on",
		"uptime-kuma":             "on",
		"webhooks":                "on",
		"whatsapp":                "on",
		"wireguard":               "on",
	}
}

func (s *EtkeccSuite) TearDownTest() {
	s.T().Helper()
	s.v.AssertExpectations(s.T())
}

func (s *EtkeccSuite) read(file string) string {
	text, err := ioutil.ReadFile("testdata/" + file)

	s.NoError(err)

	return string(text)
}

func (s *EtkeccSuite) write(file string, content string) {
	os.WriteFile("testdata/"+file, []byte(content), 0o666)
}

func (s *EtkeccSuite) rts(r io.Reader) string {
	var buf strings.Builder
	_, err := io.Copy(&buf, r)
	s.NoError(err)

	return buf.String()
}

func (s *EtkeccSuite) TestNew() {
	etkeccExt := New(s.v)

	s.IsType(&Etkecc{}, etkeccExt)
}

func (s *EtkeccSuite) TestPassword() {
	expected := "value"
	o := &order{pass: map[string]string{}}
	o.pass["test"] = "value"

	actual := o.password("test")

	s.Equal(expected, actual)
}

func (s *EtkeccSuite) TestExecute_Turnkey() {
	expectedQuestions := s.read("turnkey.questions.md")
	expectedOnboarding := s.read("turnkey.onboarding.md")
	expectedOnboardingHTML := s.read("turnkey.onboarding.html")
	expectedVars := s.read("turnkey.vars.yml")
	s.v.On("A", "example.com").Return(false).Once()
	s.v.On("CNAME", "example.com").Return(false).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("turnkey", s.turnkey)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("turnkey.questions.md", actualQuestions)
		s.write("turnkey.onboarding.md", actualOnboarding)
		s.write("turnkey.onboarding.html", actualOnboardingHTML)
		s.write("turnkey.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Turnkey_A() {
	expectedQuestions := s.read("turnkey_a.questions.md")
	expectedOnboarding := s.read("turnkey_a.onboarding.md")
	expectedOnboardingHTML := s.read("turnkey_a.onboarding.html")
	expectedVars := s.read("turnkey_a.vars.yml")
	s.v.On("A", "example.com").Return(true).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("turnkey", s.turnkey)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("turnkey_a.questions.md", actualQuestions)
		s.write("turnkey_a.onboarding.md", actualOnboarding)
		s.write("turnkey_a.onboarding.html", actualOnboardingHTML)
		s.write("turnkey_a.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Turnkey_Full() {
	expectedQuestions := s.read("turnkey_full.questions.md")
	expectedOnboarding := s.read("turnkey_full.onboarding.md")
	expectedOnboardingHTML := s.read("turnkey_full.onboarding.html")
	expectedVars := s.read("turnkey_full.vars.yml")
	s.v.On("A", "example.com").Return(false).Once()
	s.v.On("CNAME", "example.com").Return(false).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("turnkey", s.turnkeyFull)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("turnkey_full.questions.md", actualQuestions)
		s.write("turnkey_full.onboarding.md", actualOnboarding)
		s.write("turnkey_full.onboarding.html", actualOnboardingHTML)
		s.write("turnkey_full.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Turnkey_Full_A() {
	expectedQuestions := s.read("turnkey_full_a.questions.md")
	expectedOnboarding := s.read("turnkey_full_a.onboarding.md")
	expectedOnboardingHTML := s.read("turnkey_full_a.onboarding.html")
	expectedVars := s.read("turnkey_full_a.vars.yml")
	s.v.On("A", "example.com").Return(true).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("turnkey", s.turnkeyFull)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("turnkey_full_a.questions.md", actualQuestions)
		s.write("turnkey_full_a.onboarding.md", actualOnboarding)
		s.write("turnkey_full_a.onboarding.html", actualOnboardingHTML)
		s.write("turnkey_full_a.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Byos() {
	expectedQuestions := s.read("byos.questions.md")
	expectedOnboarding := s.read("byos.onboarding.md")
	expectedOnboardingHTML := s.read("byos.onboarding.html")
	expectedVars := s.read("byos.vars.yml")
	s.v.On("A", "example.com").Return(false).Once()
	s.v.On("CNAME", "example.com").Return(false).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("byos", s.byos)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("byos.questions.md", actualQuestions)
		s.write("byos.onboarding.md", actualOnboarding)
		s.write("byos.onboarding.html", actualOnboardingHTML)
		s.write("byos.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Byos_A() {
	expectedQuestions := s.read("byos_a.questions.md")
	expectedOnboarding := s.read("byos_a.onboarding.md")
	expectedOnboardingHTML := s.read("byos_a.onboarding.html")
	expectedVars := s.read("byos_a.vars.yml")
	s.v.On("A", "example.com").Return(true).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("byos", s.byos)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("byos_a.questions.md", actualQuestions)
		s.write("byos_a.onboarding.md", actualOnboarding)
		s.write("byos_a.onboarding.html", actualOnboardingHTML)
		s.write("byos_a.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Byos_Full() {
	expectedQuestions := s.read("byos_full.questions.md")
	expectedOnboarding := s.read("byos_full.onboarding.md")
	expectedOnboardingHTML := s.read("byos_full.onboarding.html")
	expectedVars := s.read("byos_full.vars.yml")
	s.v.On("A", "example.com").Return(false).Once()
	s.v.On("CNAME", "example.com").Return(false).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("byos", s.byosFull)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("byos_full.questions.md", actualQuestions)
		s.write("byos_full.onboarding.md", actualOnboarding)
		s.write("byos_full.onboarding.html", actualOnboardingHTML)
		s.write("byos_full.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func (s *EtkeccSuite) TestExecute_Byos_Full_A() {
	expectedQuestions := s.read("byos_full_a.questions.md")
	expectedOnboarding := s.read("byos_full_a.onboarding.md")
	expectedOnboardingHTML := s.read("byos_full_a.onboarding.html")
	expectedVars := s.read("byos_full_a.vars.yml")
	s.v.On("A", "example.com").Return(true).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute("byos", s.byosFull)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("byos_full_a.questions.md", actualQuestions)
		s.write("byos_full_a.onboarding.md", actualOnboarding)
		s.write("byos_full_a.onboarding.html", actualOnboardingHTML)
		s.write("byos_full_a.vars.yml", actualVars)
	}

	s.Equal(expectedQuestions, actualQuestions)
	s.Equal(expectedOnboarding, actualOnboarding)
	s.Equal(expectedOnboardingHTML, actualOnboardingHTML)
	s.Equal(expectedVars, actualVars)
}

func TestEtkeccSuite(t *testing.T) {
	suite.Run(t, new(EtkeccSuite))
}
