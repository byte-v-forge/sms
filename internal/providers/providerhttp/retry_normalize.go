package providerhttp

import "github.com/byte-v-forge/sms/internal/platform/httpx"

func normalizeRetryPolicy(policy RetryPolicy) RetryPolicy {
	if policy.Mode == RetryNone {
		policy.Attempts = 1
	}
	if policy.Attempts <= 0 {
		policy.Attempts = 1
	}
	if policy.BaseDelay <= 0 {
		policy.BaseDelay = defaultRetryBaseDelay
	}
	if policy.MaxDelay <= 0 {
		policy.MaxDelay = defaultRetryMaxDelay
	}
	if policy.MaxRetryAfter <= 0 {
		policy.MaxRetryAfter = policy.MaxDelay
	}
	if policy.MaxBodyBytes <= 0 {
		policy.MaxBodyBytes = httpx.DefaultMaxBodyBytes
	}
	return policy
}
