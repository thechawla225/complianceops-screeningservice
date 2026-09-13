package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScreen_ClearVerdictForUnrelatedNames(t *testing.T) {
	mux := newTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"debtorName":   "Acme Exports Ltd",
		"creditorName": "Example Trading Co",
	})
	req := httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var got map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got["verdict"] != "clear" {
		t.Errorf("expected verdict=clear, got %v", got["verdict"])
	}
}

func TestScreen_FlaggedVerdictForWatchlistMatch(t *testing.T) {
	mux := newTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"debtorName":   "NORTHWIND TRADING CONSORTIUM",
		"creditorName": "Acme Exports Ltd",
	})
	req := httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var got map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got["verdict"] != "flagged" {
		t.Errorf("expected verdict=flagged, got %v", got["verdict"])
	}
	if got["matchedEntity"] != "NORTHWIND TRADING CONSORTIUM" {
		t.Errorf("expected matchedEntity=NORTHWIND TRADING CONSORTIUM, got %v", got["matchedEntity"])
	}
}

func TestScreen_FlaggedVerdictViaAlias(t *testing.T) {
	mux := newTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"debtorName":   "NORTHWIND TRADING CO",
		"creditorName": "Acme Exports Ltd",
	})
	req := httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var got map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got["verdict"] != "flagged" {
		t.Errorf("expected verdict=flagged, got %v", got["verdict"])
	}
	if got["matchedEntity"] != "NORTHWIND TRADING CONSORTIUM" {
		t.Errorf("expected matchedEntity to report the canonical name, got %v", got["matchedEntity"])
	}
}

func TestScreen_RejectsMissingFields(t *testing.T) {
	mux := newTestRouter(t)
	body, _ := json.Marshal(map[string]string{"debtorName": "Acme Exports Ltd"})
	req := httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	var got map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got["error"] != "validation_error" {
		t.Errorf("expected error=validation_error, got %v", got["error"])
	}
}

func TestScreen_RejectsMalformedJSON(t *testing.T) {
	mux := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
