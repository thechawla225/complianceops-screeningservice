package api
 
import (
	"net/http"
	"${MODULE}/internal/screening"
)
 
//Added Logic in Health file to actually validate of REDIS is functional
func healthzHandler(engine *screening.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := engine.Redis.Ping(r.Context()).Err(); err != nil {
			writeError(w, http.StatusServiceUnavailable, "dependency_unavailable", "cannot reach redis")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}