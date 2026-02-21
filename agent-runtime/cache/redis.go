package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func Set(key, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
