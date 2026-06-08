package providerhttp

import "github.com/byte-v-forge/sms/internal/platform/httpclient"

func shouldRetryError(err error, attempt int, policy RetryPolicy) bool {
	return policy.Mode != RetryNone && attempt+1 < policy.Attempts && httpclient.IsRetryableTransportError(err)
}

func shouldRetryStatus(status int, attempt int, policy RetryPolicy) bool {
	return policy.Mode != RetryNone && attempt+1 < policy.Attempts && policy.RetryableStatus != nil && policy.RetryableStatus(status)
}
