package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Cache *redis.Client
var ctx = context.Background()

func Set(key string, val string) error {
	return Cache.Set(ctx, key, val, 0).Err()
}
func Get(key string) (string, error) {
	return Cache.Get(ctx, key).Result()
}

func KeyExist(key string) (int64, error) {
	return Cache.Exists(ctx, key).Result()
}
func ClearKeys() (int64, error) {
	return Cache.Del(ctx, "fixture", "weekIndex", "matchIndex").Result()
}

func Init() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	Cache = rdb
	return rdb
}
