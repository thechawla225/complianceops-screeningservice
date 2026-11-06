package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics (Phase 08.1): the same RED metrics and buckets as the Python services, so dashboards and SLO
// queries work the same everywhere, plus screening's domain metrics. Served at GET /metrics (router.go) with the Go
// runtime and process metrics the default registry already includes.

// 0.3 s is a bucket boundary on purpose: the end-to-end SLO is "95% of requests faster than 300 ms", and a histogram
// can only answer "how many were faster than X" when X is one of its bucket boundaries.
var latencyBuckets = []float64{0.01, 0.025, 0.05, 0.1, 0.2, 0.3, 0.5, 1, 2.5, 5}

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests handled, by route template and status code.",
	}, []string{"method", "route", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds, by route template.",
		Buckets: latencyBuckets,
	}, []string{"method", "route"})

	screeningVerdicts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "screening_verdicts_total",
		Help: "Screening verdicts returned, by result (clear or flagged).",
	}, []string{"result"})

	// 0 means every screening request fails closed (watchlist not loaded): alert on it instead of finding out from
	// failed payments. Updated by the health probe, which kubelet calls every 10 s.
	watchlistEntries = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "screening_watchlist_entries",
		Help: "Watchlist entries currently loaded in Redis.",
	})
)

// Probes and scrapes aren't counted, so probe traffic can't make availability look better than it is.
var excludedPaths = map[string]bool{"/healthz": true, "/metrics": true}

// statusRecorder remembers the status code a handler wrote.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Instrument wraps the router: counts and times every request by route template and status code.
func Instrument(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if excludedPaths[r.URL.Path] {
			mux.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		mux.ServeHTTP(rec, r)

		// The mux records the matched pattern ("POST /screen") on the request. Use its path part, never the raw URL:
		// raw paths would create a time series per distinct URL (unbounded label cardinality).
		route := "unmatched"
		if _, path, found := strings.Cut(r.Pattern, " "); found {
			route = path
		} else if r.Pattern != "" {
			route = r.Pattern
		}
		httpRequests.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
