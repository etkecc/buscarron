package etkecc

import (
	"encoding/json"
	"strings"
)

// hImage is "Ubuntu 22.04" image name
const hImage = "ubuntu-22.04"

var (
	// hLocations between human-readable name and actual name
	hLocations = map[string]string{
		"ashburn":     "ash",
		"falkenstein": "fsn1",
		"helsinki":    "hel1",
		"hillsboro":   "hil",
		"nuremberg":   "nbg1",
	}

	// hFirewall is the "matrix" firewall
	hFirewall = map[string]int{
		"firewall": 124003,
	}

	// hKeys is list of ssh keys names
	hKeys = []string{"first", "second", "third"}
)

type hVPSRequest struct {
	Name      string           `json:"name"`
	Size      string           `json:"server_type"`
	Image     string           `json:"image"`
	Firewalls []map[string]int `json:"firewalls"`
	SSHKeys   []string         `json:"ssh_keys"`
	Location  string           `json:"location"`
}

func (o *order) generateVPSCommand() string {
	location, ok := hLocations[strings.ToLower(o.get("turnkey-location"))]
	if !ok {
		location = "fsn1"
	}
	var size string
	sizeParts := strings.Split(o.get("turnkey"), "-")
	if len(sizeParts) < 2 {
		size = "cx11"
	} else {
		size = sizeParts[1]
	}

	req := &hVPSRequest{
		Name:      o.get("domain"),
		Size:      size,
		Image:     hImage,
		Firewalls: []map[string]int{hFirewall},
		SSHKeys:   hKeys,
		Location:  location,
	}

	return o.getVPSCurl(req)
}

func (o *order) getVPSCurl(req *hVPSRequest) string {
	reqb, _ := json.Marshal(&req) //nolint:errcheck
	reqs := strings.ReplaceAll(string(reqb), "\"", "\\\"")

	var cmd strings.Builder
	cmd.WriteString("curl -X \"POST\" \"https://api.hetzner.cloud/v1/servers\" ")
	cmd.WriteString("-H \"Content-Type: application/json\" ")
	cmd.WriteString("-H \"Authorization: Bearer $HETZNER_API_TOKEN_CLOUD\" ")
	cmd.WriteString("-d \"")
	cmd.WriteString(reqs)
	cmd.WriteString("\"\n")

	return cmd.String()
}
