package etkecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var (
	gpURL  = os.Getenv("BUSCARRON_GP_URL")
	gpUser = os.Getenv("BUSCARRON_GP_USER")
	gpPass = os.Getenv("BUSCARRON_GP_PASS")
)

type gpreq struct {
	Message string    `json:"message"`
	Files   []*gpfile `json:"files"`
}

type gpfile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Action  string `json:"action"`
	Line    string `json:"line,omitempty"`
	Regex   string `json:"regex,omitempty"`
}

func MarkAsPaid(log *zerolog.Logger, domain string) {
	if gpURL == "" || gpUser == "" || gpPass == "" {
		log.Warn().Msg("gp disabled")
	}

	req := &gpreq{
		Message: domain + " - paid",
		Files: []*gpfile{
			{
				Path:    "host_vars/" + domain + "/vars.yml",
				Action:  "append",
				Line:    "# etke services",
				Regex:   "^etke_subscription_confirmed.*",
				Content: "etke_subscription_confirmed: yes",
			},
		},
	}
	reqb, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(err).Msg("gp request marshal failed")
		return
	}
	r, err := http.NewRequest("POST", gpURL+"/post", bytes.NewReader(reqb))
	if err != nil {
		log.Error().Err(err).Msg("gp request failed")
		return
	}
	r.SetBasicAuth(gpUser, gpPass)
	r.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Error().Err(err).Msg("gp request failed")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Error().Str("status", resp.Status).Msg("gp request failed")
	}
}

func (o *order) toGP(hosts string) error {
	if gpURL == "" || gpUser == "" || gpPass == "" || o.test {
		return fmt.Errorf("disabled")
	}

	var data strings.Builder
	data.WriteString(fmt.Sprintf("- %s <%s, @%s:%s>\n", o.domain, o.get("email"), o.get("username"), o.domain))
	data.WriteString("    > Ko-Fi\n")
	data.WriteString("    * [ ] ")
	if o.hosting != "" {
		data.WriteString("TURNKEY")
	} else {
		data.WriteString("subscription")
	}
	if o.has("service-email") {
		data.WriteString(" + email")
	}

	req := &gpreq{
		Message: o.domain + " - init",
		Files:   make([]*gpfile, 0, len(o.files)),
	}
	for _, file := range o.files {
		req.Files = append(req.Files, &gpfile{
			Path:    "host_vars/" + o.domain + "/" + file.FileName,
			Action:  "create",
			Content: string(file.ContentBytes),
		})
	}
	req.Files = append(req.Files,
		&gpfile{
			Path:    "hosts",
			Action:  "append",
			Line:    "[setup]",
			Content: hosts,
		},
		&gpfile{
			Path:    "data.md",
			Action:  "prepend",
			Line:    "## Subscription",
			Content: data.String(),
		})
	reqb, err := json.Marshal(req)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", gpURL+"/post", bytes.NewReader(reqb))
	if err != nil {
		return err
	}
	r.SetBasicAuth(gpUser, gpPass)
	r.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(resp.Status)
	}
	return nil
}
