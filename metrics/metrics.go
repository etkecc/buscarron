package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	banUsers = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "buscarron_ban_users",
			Help: "The total number of banned users",
		},
		[]string{"reason", "form"})
	banRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "buscarron_ban_requests",
		Help: "The total number of requests by banned users",
	})
	submissions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "buscarron_submissions",
			Help: "The total number of successful submissions",
		},
		[]string{"form"})
)

// InitMetrics registers /metrics endpoint within default http serve mux
func InitMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

// BanRequest increments count of total requests by banned users
func BanRequest() {
	banRequests.Inc()
}

// BanUser increments counter of banned users
func BanUser(reason, form string) {
	banUsers.With(prometheus.Labels{
		"reason": reason,
		"form":   form,
	}).Inc()
}

// Submission increments count of total successful submissions
func Submission(form string) {
	submissions.With(prometheus.Labels{
		"form": form,
	}).Inc()
}
