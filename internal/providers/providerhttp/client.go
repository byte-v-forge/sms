package providerhttp

import (
	"context"
	"net/http"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestFactory func(context.Context) (*http.Request, error)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func Do(ctx context.Context, doer HTTPDoer, factory RequestFactory, policy RetryPolicy) (Response, error) {
	policy = normalizeRetryPolicy(policy)
	var lastErr error
	for attempt := 0; attempt < policy.Attempts; attempt++ {
		response, err := doOnce(ctx, doer, factory, policy.MaxBodyBytes)
		if err != nil {
			lastErr = err
			if !shouldRetryError(err, attempt, policy) {
				return Response{}, err
			}
			if err := timex.Sleep(ctx, retryDelay(attempt, policy, nil)); err != nil {
				return Response{}, err
			}
			continue
		}
		if !shouldRetryStatus(response.StatusCode, attempt, policy) {
			return response, nil
		}
		if err := timex.Sleep(ctx, retryDelay(attempt, policy, response.Header)); err != nil {
			return Response{}, err
		}
	}
	if lastErr != nil {
		return Response{}, lastErr
	}
	return Response{}, ctx.Err()
}
