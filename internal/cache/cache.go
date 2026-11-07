package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb    *redis.Client
	logger *slog.Logger
}

func New(rdb *redis.Client, logger *slog.Logger) *Cache {
	return &Cache{
		rdb:    rdb,
		logger: logger,
	}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool) {
	val, err := c.rdb.Get(ctx, key).Bytes()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			c.logger.Warn("cache GET gagal - diperlakuan sebagai miss", "key", key, "error", err)
		}
		return nil, false
	}

	return val, true
}

func (c *Cache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) {
	if err := c.rdb.Set(ctx, key, val, ttl).Err(); err != nil {
		c.logger.Warn("cache SET gagal - diabaikan", "key", key, "error", err)
	}
}
