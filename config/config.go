package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis struct {
		Host string
		Port string
	}
	Bucket struct {
		Time     string
		Expiry   string
		Treshold string
	}
}

func GetConfig() (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("error load env file")
	}
	config.Redis.Host = os.Getenv("REDIS_HOST")
	config.Redis.Port = os.Getenv("REDIS_PORT")
	config.Bucket.Time = os.Getenv("BUCKET_TIME")
	config.Bucket.Expiry = os.Getenv("BUCKET_EXPIRY")
	config.Bucket.Treshold = os.Getenv("BUCKET_THRESHOLD")
	
	return
}
