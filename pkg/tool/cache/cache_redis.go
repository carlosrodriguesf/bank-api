package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return redisCache{
		client: client,
	}
}

func (r redisCache) Set(ctx context.Context, key string, value interface{}, d time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, d).Err()
}

func (r redisCache) Get(ctx context.Context, key string, value interface{}) error {
	result := r.client.Get(ctx, key)
	if err := result.Err(); err != nil {
		return err
	}
	return json.Unmarshal([]byte(result.Val()), value)
}

func (r redisCache) GetUpdating(ctx context.Context, key string, value interface{}, d time.Duration) error {
	if err := r.Get(ctx, key, value); err != nil {
		return err
	}
	return r.Set(ctx, key, value, d)
}

func (r redisCache) Close() error {
	return r.client.Close()
}

func (r redisCache) IsErrCacheMissing(err error) bool {
	return err == redis.Nil
}
