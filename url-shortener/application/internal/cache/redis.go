package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(address string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	return &RedisCache{client: client}
}

func (r *RedisCache) Set(key string, val string) error {
	return r.client.Set(context.TODO(), key, val, time.Minute*30).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(context.TODO(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}
