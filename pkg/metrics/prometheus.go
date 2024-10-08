package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
)

var (
	routePrefixes = map[string]string{
		"/api/v1/healthcheck": "healthcheck",
		"/api/v1/auth":        "auth",
		"/api/v1/user":        "user",
		"/api/v1/films":       "films",
		"/api/v1/collections": "collections",
	}

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of the duration of HTTP requests.",
		},
		[]string{"method", "handler"},
	)

	statusCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_status_count",
			Help: "Total number of HTTP responses by status code",
		},
		[]string{"status"},
	)
)

// Register metrics
func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(statusCount)
}

func RecordRequestDuration(r *http.Request, duration float64) {
	requestDuration.WithLabelValues(r.Method, getResourceType(r)).Observe(duration)
}

func IncStatusCount(status int) {
	statusCount.WithLabelValues(fmt.Sprintf("%d", status)).Inc()
}

func getResourceType(r *http.Request) string {
	for prefix, handler := range routePrefixes {
		if strings.HasPrefix(r.RequestURI, prefix) {
			return handler
		}
	}

	if r.RequestURI == "/api" {
		return "api"
	}

	return "unknown"
}
