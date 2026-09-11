package api
 
import (
	"encoding/json"
	"net/http"
	"strings"
 
	"${MODULE}/internal/screening"
)
 

type screenRequest struct {
	DebtorName   string \`json:"debtorName"\`
	CreditorName string \`json:"creditorName"\`
}
 
func screen(w http.ResponseWriter, r *http.Request) {
	var req screenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error",
			"request body must be valid JSON matching {debtorName, creditorName}")
		return
	}
 
	if strings.TrimSpace(req.DebtorName) == "" || strings.TrimSpace(req.CreditorName) == "" {
		writeError(w, http.StatusBadRequest, "validation_error",
			"debtorName and creditorName are both required")
		return
	}
 
	result := screening.Verdict(req.DebtorName, req.CreditorName)
	writeJSON(w, http.StatusOK, result)
}