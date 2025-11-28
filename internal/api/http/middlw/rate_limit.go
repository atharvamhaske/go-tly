package middlw

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// this is a api rate limiter which means how many times user can hit and don't make ddos or spam
type ShortenLimiter struct {
	rdb    *redis.Client
	Limit  int //max requests allowed
	Window time.Duration
}

func ShortenRateLimiter(rdb *redis.Client) *ShortenLimiter {
	return &ShortenLimiter{
		rdb:    rdb,
		Limit:  20, //20 times api call
		Window: time.Minute,
	}
}

func (s *ShortenLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ip := c.RealIP()
			key := fmt.Sprintf("shorten:limit:%s", ip)
			ctx := c.Request().Context()

			count, err := s.rdb.Incr(ctx, key).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error": "rate limiter unavailable",
				})
			}
			if count == 1 {
				s.rdb.Expire(ctx, key, s.Window)
			}
			if count > int64(s.Limit) {
				ttl, _ := s.rdb.TTL(ctx, key).Result()
				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"error":      "Too many requests, please try again after some time",
					"retry_time": int(ttl.Seconds()),
				})
			}
			return next(c)
		}
	}
}
