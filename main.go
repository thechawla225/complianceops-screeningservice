package main

import (
	"${MODULE}/internal/api"
	"${MODULE}/internal/screening"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	redisAddr := getenv("REDIS_ADDR", "localhost:6379")
	watchlistPath := getenv("WATCHLIST_PATH", "watchlist.json")

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("cannot reach redis at %s: %v", redisAddr, err)
	}

	//Load the File into REDIS
	count, err := screening.SeedWatchlist(ctx, rdb, watchlistPath)
	if err != nil {
		log.Fatalf("failed to load watchlist from %s: %v", watchlistPath, err)
	}
	log.Printf("loaded %d watchlist entries into redis", count)

	engine := screening.NewEngine(rdb)
	mux := api.NewRouter(engine)

	log.Println("screening service listening on :8002")
	if err := http.ListenAndServe(":8002", mux); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
