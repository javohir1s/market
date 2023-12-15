package redis

import (
	"context"
	"fmt"
	"market_system/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	cfg    *config.Config
	client *redis.Client
}

func NewConnectionRedis(cfg *config.Config) (*Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       int(cfg.RedisDatabase),
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("pong:", pong)

	return &Cache{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *Cache) SetX(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	err := c.client.Set(ctx, key, value, expire).Err()
	return err
}

func (c *Cache) GetX(ctx context.Context, key string) ([]byte, error) {

	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}
