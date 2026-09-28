package screening

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const matchThreshold = 0.85

func (e *Engine) Verdict(ctx context.Context, debtorName, creditorName string) (Result, error) {
	entries, err := FetchWatchlist(ctx, e.Redis)
	if err != nil {
		return Result{}, fmt.Errorf("fetching watchlist: %w", err)
	}

	ref := newScreeningRef()

	if canonical, matched := matchAgainstWatchlist(debtorName, entries); matched {
		return Result{Verdict: "flagged", MatchedEntity: &canonical, ScreeningRef: ref}, nil
	}
	if canonical, matched := matchAgainstWatchlist(creditorName, entries); matched {
		return Result{Verdict: "flagged", MatchedEntity: &canonical, ScreeningRef: ref}, nil
	}

	return Result{Verdict: "clear", ScreeningRef: ref}, nil
}

func matchAgainstWatchlist(name string, entries []Entry) (string, bool) {
	if strings.TrimSpace(name) == "" {
		return "", false
	}
	for _, entry := range entries {
		if similarity(name, entry.Name) >= matchThreshold {
			return entry.Name, true
		}
		for _, alias := range entry.AlternateNames {
			if similarity(name, alias) >= matchThreshold {
				return entry.Name, true
			}
		}
	}
	return "", false
}

func newScreeningRef() string {
	return fmt.Sprintf("SCR-%d", time.Now().UnixNano())
}
