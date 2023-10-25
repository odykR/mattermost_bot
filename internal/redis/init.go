package redis

import (
	"github.com/go-redis/redis"
)

func NewClient(redisAddr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	res := rdb.Ping()

	return rdb, res.Err()
}
