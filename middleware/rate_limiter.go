package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/rate-limiter-example/config"
	"github.com/vnurhaqiqi/rate-limiter-example/infra"
)

func getKey(IP string, bucketTime int) string {
	currentMinute := time.Now().Unix() / int64(bucketTime)
	IP = fmt.Sprintf("%s:%s", IP, strconv.FormatInt(currentMinute, 10))

	return IP
}

func CustomRateLimiter(cfg config.Config, cache infra.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bucketTime, _ := strconv.Atoi(cfg.Bucket.Time)

		IPAddress := ctx.RemoteIP()
		IPAddress = getKey(IPAddress, bucketTime)

		val, err := cache.GetIP(IPAddress)
		if err != nil {
			expiry, _ := strconv.Atoi(cfg.Bucket.Expiry)
			
			err = cache.SetIP(IPAddress, infra.DurationSecond, expiry, 0)
			if err != nil {
				log.Error().Err(err)
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			if val > cfg.Bucket.Treshold {
				err := errors.New("max request reached")
				log.Warn().Err(err)
				ctx.String(http.StatusTooManyRequests, err.Error())

				return
			}
		}
		ctx.Next()
	}
}
