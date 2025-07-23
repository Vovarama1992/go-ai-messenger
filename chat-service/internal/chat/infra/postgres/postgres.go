package postgres

import (
	"database/sql"
	"time"

	"github.com/sony/gobreaker"
)

func NewPostgresConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)
	return db, nil
}

func NewBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "postgres-breaker",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     10 * time.Second,
	})
}
