package screening

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func seedEntries(t *testing.T, rdb *redis.Client, entries []Entry) {
	t.Helper()
	fields := make(map[string]interface{}, len(entries))
	for _, e := range entries {
		b, err := json.Marshal(e)
		if err != nil {
			t.Fatalf("marshal entry uid %d: %v", e.UID, err)
		}
		fields[strconv.Itoa(e.UID)] = b
	}
	if err := rdb.HSet(context.Background(), WatchlistKey, fields).Err(); err != nil {
		t.Fatalf("seeding test watchlist: %v", err)
	}
}

func newTestEngine(t *testing.T) *Engine {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	seedEntries(t, rdb, []Entry{
		{
			UID: 40347, Name: "NORTHWIND TRADING CONSORTIUM", Type: "entity",
			Programs: []string{"RUSSIA-EO14024"}, AlternateNames: []string{"NORTHWIND TRADING CO", "NTC GROUP"},
		},
		{
			UID: 40512, Name: "ZEPHYR MARITIME HOLDINGS LTD", Type: "entity",
			Programs: []string{"DPRK3"}, AlternateNames: []string{"ZEPHYR MARITIME", "ZMH LTD"},
		},
	})
	return NewEngine(rdb)
}

func TestVerdict(t *testing.T) {
	engine := newTestEngine(t)

	cases := []struct {
		name         string
		debtorName   string
		creditorName string
		wantVerdict  string
		wantMatch    string // "" means no match expected
	}{
		{"exact debtor match", "NORTHWIND TRADING CONSORTIUM", "Acme Exports Ltd", "flagged", "NORTHWIND TRADING CONSORTIUM"},
		{"exact creditor match", "Acme Exports Ltd", "ZEPHYR MARITIME HOLDINGS LTD", "flagged", "ZEPHYR MARITIME HOLDINGS LTD"},
		{"close typo on primary name still flags", "NORTHWIND TRADING CONSORTIM", "Acme Exports Ltd", "flagged", "NORTHWIND TRADING CONSORTIUM"},
		{"alias match reports the canonical name", "NORTHWIND TRADING CO", "Acme Exports Ltd", "flagged", "NORTHWIND TRADING CONSORTIUM"},
		{"unrelated names stay clear", "Acme Exports Ltd", "Example Trading Co", "clear", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := engine.Verdict(context.Background(), tc.debtorName, tc.creditorName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Verdict != tc.wantVerdict {
				t.Errorf("verdict = %q, want %q", result.Verdict, tc.wantVerdict)
			}
			if tc.wantMatch == "" {
				if result.MatchedEntity != nil {
					t.Errorf("expected no matchedEntity, got %q", *result.MatchedEntity)
				}
				if result.ScreeningRef == "" {
					t.Errorf("expected a non-empty screeningRef even when clear")
				}
			} else {
				if result.MatchedEntity == nil || *result.MatchedEntity != tc.wantMatch {
					t.Errorf("matchedEntity = %v, want %q", result.MatchedEntity, tc.wantMatch)
				}
			}
		})
	}
}
