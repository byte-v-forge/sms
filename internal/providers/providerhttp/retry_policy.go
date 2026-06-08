package providerhttp

import (
	"net/http"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/httpx"
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

func NoRetry() RetryPolicy {
	return RetryPolicy{Mode: RetryNone, Attempts: 1, MaxBodyBytes: httpx.DefaultMaxBodyBytes}
}

func DefaultRetry() RetryPolicy {
	return RetryPolicy{
		Mode:            RetryIdempotent,
		Attempts:        3,
		BaseDelay:       200 * time.Millisecond,
		MaxDelay:        2 * time.Second,
		MaxRetryAfter:   2 * time.Second,
		MaxBodyBytes:    httpx.DefaultMaxBodyBytes,
		Jitter:          100 * time.Millisecond,
		RetryableStatus: DefaultRetryableStatus,
	}
}

func DefaultRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

func normalizeRetryPolicy(policy RetryPolicy) RetryPolicy {
	if policy.Mode == RetryNone {
		policy.Attempts = 1
	}
	if policy.Attempts <= 0 {
		policy.Attempts = 1
	}
	if policy.BaseDelay <= 0 {
		policy.BaseDelay = 200 * time.Millisecond
	}
	if policy.MaxDelay <= 0 {
		policy.MaxDelay = 2 * time.Second
	}
	if policy.MaxRetryAfter <= 0 {
		policy.MaxRetryAfter = policy.MaxDelay
	}
	if policy.MaxBodyBytes <= 0 {
		policy.MaxBodyBytes = httpx.DefaultMaxBodyBytes
	}
	return policy
}
