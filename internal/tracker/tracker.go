package tracker

import "github.com/redis/go-redis/v9"

// redis client
var Rdb *redis.Client

// initializer function (runs on any import of current file) to initialize redis client
func init() {
	// initializing redis client
	Rdb = RedisInit()
}
