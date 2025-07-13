package grpcutil

import (
	"time"

	"github.com/sony/gobreaker"
)

type BreakerMode string

const (
	ModeConsecutive BreakerMode = "consecutive" // Открываем после N подряд ошибок
	ModeErrorRate   BreakerMode = "rate"        // Открываем при превышении процента ошибок за интервал времени
)

type BreakerConfig struct {
	Name string
	Mode BreakerMode

	// Для ModeConsecutive
	FailureThreshold uint32 // Сколько подряд ошибок чтобы открыть breaker

	// Для ModeErrorRate
	ErrorRateThreshold float64       // Процент ошибок от 0 до 1 (например, 0.5 = 50%)
	Interval           time.Duration // За какой период считаем процент ошибок

	// Общее для обоих режимов
	OpenStateTimeout    time.Duration // Сколько держим в состоянии OPEN
	HalfOpenMaxRequests uint32        // Сколько пробуем в HALF-OPEN (обычно 1)
}

func NewBreaker(cfg BreakerConfig) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.HalfOpenMaxRequests,
		Timeout:     cfg.OpenStateTimeout,
	}

	if cfg.Mode == ModeConsecutive {
		settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.FailureThreshold
		}
		// В режиме consecutive Interval не нужен, но gobreaker требует хоть какое-то значение.
		settings.Interval = 0
	}

	if cfg.Mode == ModeErrorRate {
		settings.Interval = cfg.Interval
		settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
			total := counts.Requests
			if total == 0 {
				return false
			}
			errorRate := float64(counts.TotalFailures) / float64(total)
			return errorRate >= cfg.ErrorRateThreshold
		}
	}

	return gobreaker.NewCircuitBreaker(settings)
}
