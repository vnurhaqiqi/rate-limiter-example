package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/rate-limiter-example/config"
	"github.com/vnurhaqiqi/rate-limiter-example/infra"
	"github.com/vnurhaqiqi/rate-limiter-example/middleware"
)

func main() {
	r := gin.Default()

	// get all configs
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	// provide cache client
	cache := infra.ProvideClient(cfg)

	// rate limiter section
	customeRateLimiter := middleware.CustomRateLimiter(cfg, *cache)

	r.GET("/custom-rate-limiter", customeRateLimiter, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})
}
