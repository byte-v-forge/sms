package providerhttp

import (
	"time"
)

type RetryMode int

const (
	RetryNone RetryMode = iota
	RetryIdempotent
)

type RetryPolicy struct {
	Mode            RetryMode
	Attempts        int
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	MaxRetryAfter   time.Duration
	MaxBodyBytes    int64
	Jitter          time.Duration
	RetryableStatus func(int) bool
}
