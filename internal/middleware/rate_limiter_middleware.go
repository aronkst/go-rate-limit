package middleware

import (
	"net"
	"net/http"

	"github.com/aronkst/go-rate-limit/internal/ratelimit"
)

func RateLimiterMiddleware(limiter ratelimit.RateLimitInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var identifier string

			token := r.Header.Get("API_KEY")
			if token != "" {
				identifier = token
			} else {
				identifier = getIPFromRemoteAddr(r.RemoteAddr)
			}

			if exceeded, err := limiter.IsLimitExceeded(identifier); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			} else if exceeded {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIPFromRemoteAddr(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
