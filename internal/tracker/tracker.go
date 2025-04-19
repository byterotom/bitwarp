package tracker

import "github.com/redis/go-redis/v9"

// redis client
var Rdb *redis.Client

func init() {
	// initializing redis client
	Rdb = RedisInit()
}
