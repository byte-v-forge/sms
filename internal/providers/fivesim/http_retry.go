package fivesim

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func fiveSimRetryPolicy(path string) providerhttp.RetryPolicy {
	var policy providerhttp.RetryPolicy
	if fiveSimRequestRetryable(path) {
		policy = providerhttp.DefaultRetry()
	} else {
		policy = providerhttp.NoRetry()
	}
	policy.MaxBodyBytes = 1 << 20
	return policy
}

func fiveSimRequestRetryable(path string) bool {
	return strings.HasPrefix(path, "/v1/guest/") || strings.HasPrefix(path, "/v1/user/check/") || path == "/v1/user/profile"
}
