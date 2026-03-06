package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
	ctx context.Context
}

func New(redisURL string) (*Cache, error) {
	var opts *redis.Options
	parsed, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{Addr: redisURL}
	} else {
		opts = parsed
	}

	rdb := redis.NewClient(opts)
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &Cache{rdb: rdb, ctx: ctx}, nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	val, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (c *Cache) Set(key string, data []byte, ttl time.Duration) {
	_ = c.rdb.Set(c.ctx, key, data, ttl).Err()
}

func (c *Cache) Close() error {
	return c.rdb.Close()
}
