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

type EtkeccSuite struct {
	suite.Suite
	v           *mocks.Validator
	ext         *Etkecc
	save        bool
	byos        map[string]string
	byosSub     map[string]string
	turnkey     map[string]string
	byosFull    map[string]string
	turnkeyFull map[string]string
}

func (s *EtkeccSuite) SetupTest() {
	s.T().Helper()
	s.v = &mocks.Validator{}
	s.ext = New(nil)
	s.ext.test = true
	s.save = false

	s.byos = map[string]string{
		"homeserver":    "synapse",
		"domain":        "https://matrix.ExAmPlE.com ",
		"username":      " tEsT ",
		"email":         "tEsT@TEST.cOm",
		"type":          "byos",
		"synapse-admin": "on",
		"etherpad":      "on",
		"smtp-relay":    "on",
	}
	s.turnkey = map[string]string{
		"homeserver": "synapse",
		"domain":     "https://matrix.ExAmPlE.com ",
		"username":   " tEsT ",
		"email":      "tEsT@TEST.cOm",
		"type":       "turnkey",
		"lang":       "wrong",
	}
	s.byosSub = map[string]string{
		"domain":        "https://higenjitsuteki.etke.host",
		"domain-type":   "subdomain",
		"username":      " tEsT ",
		"email":         "tEsT@TEST.cOm",
		"type":          "byos",
		"synapse-admin": "on",
	}
	s.byosFull = map[string]string{
		"homeserver":          "synapse",
		"domain":              "https://matrix.ExAmPlE.com ",
		"username":            " tEsT ",
		"email":               "tEsT@TEST.cOm",
		"name":                "Test",
		"notes":               "ALL, I.WANT.THEM.ALLLLLLL. Yes, on shitty 1vCPU 1GB RAM VPS.",
		"service-setup":       "on",
		"service-maintenance": "on",
		"service-email":       "on",
		"lang":                "de",
		"borg":                "on",
		"bridges-encryption":  "on",
		"buscarron":           "on",
		"cinny":               "on",
		"dimension":           "on",
		"discord":             "on",
		"element-web":         "on",
		"etherpad":            "on",
		"facebook":            "on",
		"gmessages":           "on",
		"googlechat":          "on",
		"groupme":             "on",
		"hangouts":            "on",
		"heisenbridge":        "on",
		"honoroit":            "on",
		"hydrogen":            "on",
		"instagram":           "on",
		"irc":                 "on",
		"jitsi":               "on",
		"linkedin":            "on",
		"nginx-proxy-website": "on",
		"ntfy":                "on",
		"postmoogle":          "on",
		"reminder-bot":        "on",
		"signal":              "on",
		"skype":               "on",
		"slack":               "on",
		"smtp-relay":          "on",
		"sso":                 "on",
		"stats":               "on",
		"steam":               "on",
		"sygnal":              "on",
		"telegram":            "on",
		"twitter":             "on",
		"type":                "byos",
		"webhooks":            "on",
		"whatsapp":            "on",
	}
	s.turnkeyFull = map[string]string{
		"homeserver":          "synapse",
		"domain":              "https://higenjitsuteki.etke.host",
		"domain-type":         "subdomain",
		"username":            " tEsT ",
		"email":               "tEsT@TEST.cOm",
		"name":                "Test",
		"notes":               "ALL, I.WANT.THEM.ALLLLLLL. Yes, on shitty 1vCPU 1GB RAM VPS.",
		"service-setup":       "on",
		"service-maintenance": "on",
		"service-email":       "on",
		"lang":                "invalid",
		"type":                "turnkey",
		"borg":                "on",
		"buscarron":           "on",
		"cinny":               "on",
		"dimension":           "on",
		"discord":             "on",
		"dnsmasq":             "on",
		"element-web":         "on",
		"etherpad":            "on",
		"facebook":            "on",
		"googlechat":          "on",
		"groupme":             "on",
		"hangouts":            "on",
		"heisenbridge":        "on",
		"honoroit":            "on",
		"hydrogen":            "on",
		"instagram":           "on",
		"irc":                 "on",
		"jitsi":               "on",
		"linkedin":            "on",
		"nginx-proxy-website": "on",
		"ntfy":                "on",
		"reminder-bot":        "on",
		"signal":              "on",
		"skype":               "on",
		"slack":               "on",
		"smtp-relay":          "on",
		"sso":                 "on",
		"stats":               "on",
		"steam":               "on",
		"sygnal":              "on",
		"telegram":            "on",
		"twitter":             "on",
		"webhooks":            "on",
		"whatsapp":            "on",
	}
}

func (s *EtkeccSuite) TearDownTest() {
	s.T().Helper()
	s.v.AssertExpectations(s.T())
}

func (s *EtkeccSuite) read(file string) string {
	text, err := os.ReadFile("testdata/" + file)

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

func (s *EtkeccSuite) TestExecute_Turnkey() {
	expectedQuestions := s.read("turnkey.questions.md")
	expectedOnboarding := s.read("turnkey.onboarding.md")
	expectedOnboardingHTML := s.read("turnkey.onboarding.html")
	expectedVars := s.read("turnkey.vars.yml")
	s.v.On("A", "example.com").Return(false).Once()
	s.v.On("CNAME", "example.com").Return(false).Once()
	s.v.On("GetBase", "https://matrix.example.com").Return("example.com").Once()

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "turnkey"}, s.turnkey)
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

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "turnkey"}, s.turnkey)
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
	s.v.On("A", "higenjitsuteki.etke.host").Return(false).Once()
	s.v.On("CNAME", "higenjitsuteki.etke.host").Return(false).Once()
	s.v.On("GetBase", "https://higenjitsuteki.etke.host").Return("higenjitsuteki.etke.host").Once()

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "turnkey"}, s.turnkeyFull)
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
	s.v.On("A", "higenjitsuteki.etke.host").Return(true).Once()
	s.v.On("GetBase", "https://higenjitsuteki.etke.host").Return("higenjitsuteki.etke.host").Once()

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "turnkey"}, s.turnkeyFull)
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

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "byos"}, s.byos)
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

func (s *EtkeccSuite) TestExecute_Byos_Sub() {
	expectedQuestions := s.read("byos_sub.questions.md")
	expectedOnboarding := s.read("byos_sub.onboarding.md")
	expectedOnboardingHTML := s.read("byos_sub.onboarding.html")
	expectedVars := s.read("byos_sub.vars.yml")
	s.v.On("A", "higenjitsuteki.etke.host").Return(false).Once()
	s.v.On("CNAME", "higenjitsuteki.etke.host").Return(false).Once()
	s.v.On("GetBase", "https://higenjitsuteki.etke.host").Return("higenjitsuteki.etke.host").Once()

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "byos"}, s.byosSub)
	actualVars := s.rts(files[0].Content)
	actualOnboarding := s.rts(files[1].Content)
	actualOnboardingHTML := s.rts(files[2].Content)
	// to generate output
	if s.save {
		s.write("byos_sub.questions.md", actualQuestions)
		s.write("byos_sub.onboarding.md", actualOnboarding)
		s.write("byos_sub.onboarding.html", actualOnboardingHTML)
		s.write("byos_sub.vars.yml", actualVars)
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

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "byos"}, s.byos)
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

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "byos"}, s.byosFull)
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

	actualQuestions, files := s.ext.Execute(s.v, &config.Form{Name: "byos"}, s.byosFull)
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
