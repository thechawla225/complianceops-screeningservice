package screening

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Fetch all the enties from REDIS for the screening request
func FetchWatchlist(ctx context.Context, rdb *redis.Client) ([]Entry, error) {
	raw, err := rdb.HGetAll(ctx, WatchlistKey).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uid, encoded := range raw {
		var e Entry
		if err := json.Unmarshal([]byte(encoded), &e); err != nil {
			return nil, fmt.Errorf("decoding cached watchlist entry uid %s: %w", uid, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
