package api

import (
	"complianceops-screening/internal/screening"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Defining the Endpoitns for this service
func NewRouter(engine *screening.Engine) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(engine))
	mux.HandleFunc("POST /screen", screenHandler(engine))
	mux.Handle("GET /metrics", promhttp.Handler())
	return mux
}
