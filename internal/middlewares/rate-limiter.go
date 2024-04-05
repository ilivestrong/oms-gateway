package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(limit int, duration time.Duration) func(http.Handler) http.Handler {
	limiter := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("rate limiter called")
			if limiter.TakeAvailable(1) == 0 {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
