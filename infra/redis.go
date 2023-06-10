package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vnurhaqiqi/rate-limiter-example/config"
)

var (
	DurationSecond = "second"
	DurationMinute = "minute"
	DurationHour   = "hour"
)

type Cache struct {
	Redis *redis.Client
	Ctx   context.Context
}

func RedisClient(cfg config.Config) *Cache {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       0,
	})

	return &Cache{
		Ctx:   ctx,
		Redis: redisClient,
	}
}

func (c *Cache) SetIP(key string, durationType string, expiry int, val interface{}) error {
	duration := getDurationType(durationType)

	err := c.Redis.Set(c.Ctx, key, val, time.Duration(expiry)*duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetIP(IP string) (string, error) {
	val, err := c.Redis.Get(c.Ctx, IP).Result()
	if err != nil {
		return val, err
	}

	return val, nil
}

// getDurationType this function to resolve duration type for expiry
func getDurationType(durationType string) time.Duration {
	switch durationType {
	case DurationSecond:
		return time.Second
	case DurationMinute:
		return time.Minute
	case DurationHour:
		return time.Hour
	default:
		return time.Second
	}
}
