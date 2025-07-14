package config

import (
	"os"
	"strconv"
	"time"
)

type PostgresBreakerConfig struct {
	OpenTimeout      time.Duration
	FailureThreshold uint32
	MaxRequests      uint32
}

func LoadPostgresBreakerConfig() PostgresBreakerConfig {
	openTimeoutSec, _ := strconv.Atoi(os.Getenv("CB_POSTGRES_OPEN_TIMEOUT"))
	if openTimeoutSec == 0 {
		openTimeoutSec = 10
	}

	failureThreshold, _ := strconv.Atoi(os.Getenv("CB_POSTGRES_FAILURE_THRESHOLD"))
	if failureThreshold == 0 {
		failureThreshold = 5
	}

	maxRequests, _ := strconv.Atoi(os.Getenv("CB_POSTGRES_MAX_REQUESTS"))
	if maxRequests == 0 {
		maxRequests = 1
	}

	return PostgresBreakerConfig{
		OpenTimeout:      time.Duration(openTimeoutSec) * time.Second,
		FailureThreshold: uint32(failureThreshold),
		MaxRequests:      uint32(maxRequests),
	}
}
