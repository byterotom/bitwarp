// redis.go
package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// time to live for each entry
const TTL = 10 * time.Second

// function to initialize redis client to interact with redis
func RedisInit() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return rdb
}

// function to remove expired members
func RemoveExpired(key string, ctx context.Context) {
	now := float64(time.Now().Unix())
	Rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", now))
}
