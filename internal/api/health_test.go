package api

import (
	"${MODULE}/internal/screening"
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
		UID: 40347, Name: "RANDOM RUSSIAN COMPANY", Type: "entity",
		Programs:       []string{"RUSSIA-EO14024"},
		AlternateNames: []string{"RANDOM RUSSIAN CO", "RRC GROUP"},
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
