package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewRedis(host string, port int) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})
	return &Cache{client: rdb}
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("cache miss: %s", key)
	}
	if err != nil {
		return fmt.Errorf("get from cache: %w", err)
	}
	return json.Unmarshal(val, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
