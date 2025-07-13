package httputil

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func NewRateLimiter(maxRequests int, per time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(maxRequests, per)
}
