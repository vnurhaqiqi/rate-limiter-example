package infra

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/vnurhaqiqi/rate-limiter-example/config"
)

func RedisClient(cfg config.Config) redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       0,
	})

	return *client
}
