package api
 
import "net/http"
 
//Will actually validate Health once redis is configured
func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}