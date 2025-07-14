package config

import (
	"os"
	"strconv"
	"time"

	"github.com/Vovarama1992/go-utils/pgutil"
)

func LoadDBConfig() pgutil.DBPoolConfig {
	return pgutil.DBPoolConfig{
		MaxConnLifetime:   parseDuration("USER_DB_CONN_LIFETIME", "5m"),
		MaxConnIdleTime:   parseDuration("USER_DB_CONN_IDLE", "2m"),
		HealthCheckPeriod: parseDuration("USER_DB_HEALTHCHECK_PERIOD", "1m"),
		ConnectTimeout:    5 * time.Second,
		MaxConns:          parseInt32("USER_DB_MAX_CONNS", 10),
	}
}

func parseDuration(envKey, fallback string) time.Duration {
	val := os.Getenv(envKey)
	if val == "" {
		val = fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0
	}
	return d
}

func parseInt32(envKey string, fallback int32) int32 {
	val := os.Getenv(envKey)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return int32(n)
}
