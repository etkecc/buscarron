package etkecc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/etkecc/buscarron/internal/utils"
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

func (o *order) toGP(ctx context.Context, hosts string) error {
	ctx = context.WithoutCancel(ctx)
	log := o.logger(ctx)
	if gpURL == "" || gpUser == "" || gpPass == "" || o.test {
		return fmt.Errorf("disabled")
	}

	log.Info().Msg("sending to GP")
	defer log.Info().Msg("sent to GP")
	span := utils.StartSpan(ctx, "sub.ext.etkecc.toGP")
	defer span.Finish()

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
		log.Error().Err(err).Msg("failed to marshal request")
		return err
	}
	r, err := http.NewRequestWithContext(span.Context(), http.MethodPost, gpURL+"/post", bytes.NewReader(reqb))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return err
	}
	r.SetBasicAuth(gpUser, gpPass)
	r.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Error().Err(err).Msg("failed to send request")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Error().Str("status", resp.Status).Msg("failed to send request")
		return errors.New(resp.Status)
	}
	return nil
}
