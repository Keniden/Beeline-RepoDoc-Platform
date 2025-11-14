package redisclient

import (
    "context"
    "fmt"
    "time"

    "github.com/beeline/repodoc/configs"
    "github.com/go-redis/redis/v8"
)

func New(ctx context.Context, cfg configs.StorageConfig) (*redis.Client, error) {
    opt := &redis.Options{
        Addr: cfg.RedisAddr,
    }
    client := redis.NewClient(opt)
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis ping: %w", err)
    }
    return client, nil
}

func CacheKey(repoID, key string) string {
    return fmt.Sprintf("repodoc:%s:%s", repoID, key)
}

func WithTTL(ctx context.Context, client *redis.Client, key string, value interface{}) error {
    return client.Set(ctx, key, value, 2*time.Minute).Err()
}
