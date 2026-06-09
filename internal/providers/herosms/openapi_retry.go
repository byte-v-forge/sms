package herosms

import "github.com/byte-v-forge/sms/internal/providers/providerhttp"

func heroSMSOpenAPIRetryPolicy() providerhttp.RetryPolicy {
	policy := providerhttp.DefaultRetry()
	policy.MaxBodyBytes = 8 << 20
	return policy
}

func heroSMSOpenAPINonRetryPolicy() providerhttp.RetryPolicy {
	policy := providerhttp.NoRetry()
	policy.MaxBodyBytes = 1 << 20
	return policy
}
