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
	bucketTime, _ := strconv.Atoi(cfg.Bucket.Time)
	bucketTreshold, _ := strconv.Atoi(cfg.Bucket.Treshold)
	expiry, _ := strconv.Atoi(cfg.Bucket.Expiry)

	return func(ctx *gin.Context) {
		IPAddress := ctx.Request.RemoteAddr
		IPAddress = getKey(IPAddress, bucketTime)

		// Get IP Address number of count
		counter, err := cache.GetIncrement(IPAddress)
		if err != nil {

			// Set IP Address to store in cache
			err = cache.SetKey(IPAddress, infra.DurationSecond, expiry, 0)
			if err != nil {
				handleError(ctx, err, http.StatusBadRequest)
				return
			}

			// Set increment IP Address as counter for request
			err = cache.SetIncrement(IPAddress)
			if err != nil {
				handleError(ctx, err, http.StatusBadRequest)
				return
			}
		} else {
			if counter > bucketTreshold {
				err := errors.New("max request reached")
				handleError(ctx, err, http.StatusTooManyRequests)
				return
			}
			err = cache.SetIncrement(IPAddress)
			if err != nil {
				handleError(ctx, err, http.StatusBadRequest)
				return
			}
		}
		ctx.Next()
	}
}

// error handler
func handleError(ctx *gin.Context, err error, code int) {
	log.Error().Err(err)
	ctx.String(code, err.Error())
	ctx.Abort()
}
