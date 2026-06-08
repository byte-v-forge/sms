package handlerapi

import "github.com/byte-v-forge/sms/internal/providers/providerhttp"

func handlerAPIRetryPolicy(action string) providerhttp.RetryPolicy {
	switch action {
	case "getNumberV2", "setStatus":
		policy := providerhttp.NoRetry()
		policy.MaxBodyBytes = 1 << 20
		return policy
	default:
		policy := providerhttp.DefaultRetry()
		policy.MaxBodyBytes = 1 << 20
		return policy
	}
}
