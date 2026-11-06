package api

import (
	"complianceops-screening/internal/screening"
	"net/http"
)

// Added Logic in Health file to actually validate of REDIS is functional
func healthzHandler(engine *screening.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// HLEN: 0 for a missing key, like EXISTS, and also gives the entry count for the metric
		n, err := engine.Redis.HLen(r.Context(), screening.WatchlistKey).Result()
		if err != nil {
			writeError(w, http.StatusServiceUnavailable, "dependency_unavailable", "cannot reach redis")
			return
		}
		watchlistEntries.Set(float64(n))
		// Redis came back empty: fail the probe so the pod restarts and main.go reseeds the watchlist
		if n == 0 {
			writeError(w, http.StatusServiceUnavailable, "dependency_unavailable", "watchlist not loaded")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
