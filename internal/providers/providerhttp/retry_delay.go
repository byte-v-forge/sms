package providerhttp

import (
	"net/http"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

func retryDelay(attempt int, policy RetryPolicy, header http.Header) time.Duration {
	if header != nil {
		if delay := httpx.RetryAfterMax(header, policy.MaxRetryAfter); delay > 0 {
			return delay
		}
	}
	delay := policy.BaseDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
	}
	if delay > policy.MaxDelay {
		delay = policy.MaxDelay
	}
	return delay + randomJitter(policy.Jitter)
}
