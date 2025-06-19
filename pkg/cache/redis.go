package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	Addr     string
	Password string
	DB       int
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(options Options) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &RedisCache{
		client: rdb,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
