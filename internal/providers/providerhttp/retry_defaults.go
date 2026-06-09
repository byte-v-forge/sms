package providerhttp

import (
	"net/http"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

const (
	defaultRetryAttempts      = 3
	defaultRetryBaseDelay     = 200 * time.Millisecond
	defaultRetryMaxDelay      = 2 * time.Second
	defaultRetryMaxRetryAfter = 2 * time.Second
	defaultRetryJitter        = 100 * time.Millisecond
)

func NoRetry() RetryPolicy {
	return RetryPolicy{
		Mode:         RetryNone,
		Attempts:     1,
		MaxBodyBytes: httpx.DefaultMaxBodyBytes,
	}
}

func DefaultRetry() RetryPolicy {
	return RetryPolicy{
		Mode:            RetryIdempotent,
		Attempts:        defaultRetryAttempts,
		BaseDelay:       defaultRetryBaseDelay,
		MaxDelay:        defaultRetryMaxDelay,
		MaxRetryAfter:   defaultRetryMaxRetryAfter,
		MaxBodyBytes:    httpx.DefaultMaxBodyBytes,
		Jitter:          defaultRetryJitter,
		RetryableStatus: DefaultRetryableStatus,
	}
}

func DefaultRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}
