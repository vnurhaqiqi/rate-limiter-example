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

func ProvideClient(cfg config.Config) *Cache {
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

func (c *Cache) SetKey(key string, durationType string, expiry int, val interface{}) error {
	duration := getDurationType(durationType)

	err := c.Redis.Set(c.Ctx, key, val, time.Duration(expiry)*duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetKey(key string) (string, error) {
	val, err := c.Redis.Get(c.Ctx, key).Result()
	if err != nil {
		return val, err
	}

	return val, nil
}

func (c *Cache) SetIncrement(key string) error {
	err := c.Redis.Incr(c.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetIncrement(key string) (counter int, err error) {
	counter, err = c.Redis.Get(c.Ctx, key).Int()
	if err != nil {
		return
	}
	return
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
