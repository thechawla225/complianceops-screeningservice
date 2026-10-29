package api

import (
	"complianceops-screening/internal/screening"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestRouter(t *testing.T) *http.ServeMux {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	entry := screening.Entry{
		UID: 40347, Name: "NORTHWIND TRADING CONSORTIUM", Type: "entity",
		Programs:       []string{"RUSSIA-EO14024"},
		AlternateNames: []string{"NORTHWIND TRADING CO", "NTC GROUP"},
	}
	encoded, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal test entry: %v", err)
	}
	if err := rdb.HSet(context.Background(), screening.WatchlistKey, "40347", encoded).Err(); err != nil {
		t.Fatalf("seeding test watchlist: %v", err)
	}
	return NewRouter(screening.NewEngine(rdb))
}

func TestHealthz_OkWhenRedisReachable(t *testing.T) {
	mux := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHealthz_UnavailableWhenWatchlistMissing(t *testing.T) {
	mr := miniredis.RunT(t)
	mux := NewRouter(screening.NewEngine(redis.NewClient(&redis.Options{Addr: mr.Addr()})))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}
