package api

import (
	"complianceops-screening/internal/screening"
	"net/http"
)

// Defining the Endpoitns for this service
func NewRouter(engine *screening.Engine) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(engine))
	mux.HandleFunc("POST /screen", screenHandler(engine))
	return mux
}
