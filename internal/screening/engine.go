package screening

import "github.com/redis/go-redis/v9"

// Engine holds the shared dependencies the screening service needs —
// currently just the Redis client used for the watchlist and health checks.
type Engine struct {
	Redis *redis.Client
}

func NewEngine(rdb *redis.Client) *Engine {
	return &Engine{Redis: rdb}
}
