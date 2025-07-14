package circuitbreaker

import (
	"os"
	"strconv"
	"time"

	"github.com/sony/gobreaker"
)

func NewUserServiceBreaker() *gobreaker.CircuitBreaker {
	openTimeoutSec, _ := strconv.Atoi(os.Getenv("CB_USER_OPEN_TIMEOUT"))
	if openTimeoutSec == 0 {
		openTimeoutSec = 30
	}

	failureThreshold, _ := strconv.Atoi(os.Getenv("CB_USER_FAILURE_THRESHOLD"))
	if failureThreshold == 0 {
		failureThreshold = 5
	}

	maxRequests, _ := strconv.Atoi(os.Getenv("CB_USER_MAX_REQUESTS"))
	if maxRequests == 0 {
		maxRequests = 1
	}

	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "user-service",
		MaxRequests: uint32(maxRequests),
		Interval:    time.Minute,
		Timeout:     time.Duration(openTimeoutSec) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= uint32(failureThreshold)
		},
	})
}
