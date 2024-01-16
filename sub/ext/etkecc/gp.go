package etkecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func MarkAsPaid(log *zerolog.Logger, domain, baseDomain, amount string) {
	if gpURL == "" || gpUser == "" || gpPass == "" {
		log.Warn().Msg("gp disabled")
	}

	msgdomain := domain
	if domain != baseDomain {
		msgdomain += " (or " + baseDomain + ")"
	}

	req := &gpreq{
		Message: msgdomain + " - paid",
		Files: []*gpfile{
			{
				Path:    "host_vars/" + domain + "/vars.yml",
				Action:  "replace",
				Line:    "# etke services",
				Regex:   "etke_subscription_confirmed: no",
				Content: "etke_subscription_confirmed: yes",
			},
			{
				Path:    "host_vars/" + domain + "/vars.yml",
				Action:  "append",
				Line:    "# etke services",
				Content: "etke_subscription_first_payment: " + amount,
			},
		},
	}
	if domain != baseDomain {
		req.Files = append(req.Files,
			&gpfile{
				Path:    "host_vars/" + baseDomain + "/vars.yml",
				Action:  "replace",
				Line:    "# etke services",
				Regex:   "etke_subscription_confirmed: no",
				Content: "etke_subscription_confirmed: yes",
			},
			&gpfile{
				Path:    "host_vars/" + baseDomain + "/vars.yml",
				Action:  "append",
				Line:    "# etke services",
				Content: "etke_subscription_first_payment: " + amount,
			})
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
