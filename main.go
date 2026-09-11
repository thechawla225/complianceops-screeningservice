package main
 
import (
	"log"
	"net/http"
 
	"${MODULE}/internal/api"
)
 
func main() {
	mux := api.NewRouter()
 
	log.Println("screening service listening on :8002")
	if err := http.ListenAndServe(":8002", mux); err != nil {
		log.Fatal(err)
	}
}