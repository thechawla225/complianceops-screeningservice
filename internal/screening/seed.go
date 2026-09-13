package screening

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Function to write the watchlist into REDIS
func SeedWatchlist(ctx context.Context, rdb *redis.Client, path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("reading watchlist file: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return 0, fmt.Errorf("parsing watchlist file: %w", err)
	}

	//This is to clear any existing watchlist entries before we start writing new entries
	if err := rdb.Del(ctx, WatchlistKey).Err(); err != nil {
		return 0, fmt.Errorf("clearing existing watchlist: %w", err)
	}

	//if no entries, exit
	if len(entries) == 0 {
		return 0, nil
	}

	fields := make(map[string]interface{}, len(entries))
	for _, e := range entries {
		encoded, err := json.Marshal(e)
		if err != nil {
			return 0, fmt.Errorf("encoding watchlist entry uid %d: %w", e.UID, err)
		}
		fields[strconv.Itoa(e.UID)] = encoded
	}
	if err := rdb.HSet(ctx, WatchlistKey, fields).Err(); err != nil {
		return 0, fmt.Errorf("loading watchlist into redis: %w", err)
	}

	return len(entries), nil
}
