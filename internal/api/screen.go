package api

import (
	"complianceops-screening/internal/screening"
	"encoding/json"
	"net/http"
	"strings"
)

type screenRequest struct {
	DebtorName   string `json:"debtorName"`
	CreditorName string `json:"creditorName"`
}

func screenHandler(engine *screening.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		result, err := engine.Verdict(r.Context(), req.DebtorName, req.CreditorName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "screening_error", "failed to screen against watchlist")
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}
