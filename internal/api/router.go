package api
 
import "net/http"

//defining the two routers we will be using in this service
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("POST /screen", screen)
	return mux
}