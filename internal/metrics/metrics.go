package metrics

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

var banRequests = metrics.NewCounter("buscarron_ban_requests")

// Handler for metrics
type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	metrics.WritePrometheus(w, false)
}

// BanRequest increments count of total requests by banned users
func BanRequest() {
	banRequests.Inc()
}

// BanUser increments counter of banned users
func BanUser(reason, form string) {
	metrics.GetOrCreateCounter(fmt.Sprintf("buscarron_ban_users{form=%q,reason=%q}", form, reason)).Inc()
}

// Submission increments count of total successful submissions
func Submission(form string) {
	metrics.GetOrCreateCounter(fmt.Sprintf("buscarron_submissions{form=%q}", form)).Inc()
}
