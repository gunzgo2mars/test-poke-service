package cache

import "github.com/redis/go-redis/v9"

func NewRedis(opt *redis.Options) *redis.Client {
	return redis.NewClient(opt)
}
