package etkecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
}

func (o *order) toGP() error {
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
