package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func screenOnce(t *testing.T, h http.Handler, debtor, creditor string) int {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"debtorName": debtor, "creditorName": creditor})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/screen", bytes.NewReader(body)))
	return rec.Code
}

func scrape(t *testing.T, h http.Handler) string {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/metrics: expected 200, got %d", rec.Code)
	}
	return rec.Body.String()
}

func TestMetrics_RequestsCountedByRouteTemplateAndStatus(t *testing.T) {
	h := Instrument(newTestRouter(t))
	ok := httpRequests.WithLabelValues("POST", "/screen", "200")
	bad := httpRequests.WithLabelValues("POST", "/screen", "400")
	okBefore, badBefore := testutil.ToFloat64(ok), testutil.ToFloat64(bad)

	screenOnce(t, h, "Acme Exports Ltd", "Example Trading Co")
	screenOnce(t, h, "", "")

	if got := testutil.ToFloat64(ok) - okBefore; got != 1 {
		t.Errorf("200s on /screen: expected +1, got +%v", got)
	}
	if got := testutil.ToFloat64(bad) - badBefore; got != 1 {
		t.Errorf("400s on /screen: expected +1, got +%v", got)
	}
}

func TestMetrics_UnknownPathsShareOneLabel(t *testing.T) {
	h := Instrument(newTestRouter(t))
	c := httpRequests.WithLabelValues("GET", "unmatched", "404")
	before := testutil.ToFloat64(c)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/no-such-path-123", nil))

	if got := testutil.ToFloat64(c) - before; got != 1 {
		t.Errorf("unmatched 404: expected +1, got +%v", got)
	}
	if strings.Contains(scrape(t, h), "no-such-path-123") {
		t.Error("raw path leaked into metric labels")
	}
}

func TestMetrics_VerdictsCounted(t *testing.T) {
	h := Instrument(newTestRouter(t))
	clear, flagged := screeningVerdicts.WithLabelValues("clear"), screeningVerdicts.WithLabelValues("flagged")
	clearBefore, flaggedBefore := testutil.ToFloat64(clear), testutil.ToFloat64(flagged)

	screenOnce(t, h, "Acme Exports Ltd", "Example Trading Co")
	screenOnce(t, h, "NORTHWIND TRADING CONSORTIUM", "Acme Exports Ltd")

	if got := testutil.ToFloat64(clear) - clearBefore; got != 1 {
		t.Errorf("clear verdicts: expected +1, got +%v", got)
	}
	if got := testutil.ToFloat64(flagged) - flaggedBefore; got != 1 {
		t.Errorf("flagged verdicts: expected +1, got +%v", got)
	}
}

func TestMetrics_WatchlistGaugeFollowsHealthProbe(t *testing.T) {
	h := Instrument(newTestRouter(t)) // test watchlist has 1 entry
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz: expected 200, got %d", rec.Code)
	}
	if got := testutil.ToFloat64(watchlistEntries); got != 1 {
		t.Errorf("screening_watchlist_entries: expected 1, got %v", got)
	}
}

func TestMetrics_EndpointExcludesProbesAndHasSLOBucket(t *testing.T) {
	h := Instrument(newTestRouter(t))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	screenOnce(t, h, "Acme Exports Ltd", "Example Trading Co")

	body := scrape(t, h)
	if strings.Contains(body, `route="/healthz"`) || strings.Contains(body, `route="/metrics"`) {
		t.Error("probes or scrapes were counted")
	}
	if !strings.Contains(body, `le="0.3"`) {
		t.Error(`missing the le="0.3" bucket the SLO needs`)
	}
	if !strings.Contains(body, "go_goroutines") {
		t.Error("expected Go runtime metrics from the default registry")
	}
}
