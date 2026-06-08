package providerhttp

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/httpclient"
	"github.com/byte-v-forge/sms/internal/platform/httpx"
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

func doOnce(ctx context.Context, doer HTTPDoer, factory RequestFactory, maxBodyBytes int64) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	req, err := factory(ctx)
	if err != nil {
		return Response{}, err
	}
	resp, err := doer.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := httpx.ReadLimited(resp.Body, maxBodyBytes)
	if err != nil {
		return Response{StatusCode: resp.StatusCode, Header: resp.Header}, err
	}
	return Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: body}, nil
}

func shouldRetryError(err error, attempt int, policy RetryPolicy) bool {
	return policy.Mode != RetryNone && attempt+1 < policy.Attempts && httpclient.IsRetryableTransportError(err)
}

func shouldRetryStatus(status int, attempt int, policy RetryPolicy) bool {
	return policy.Mode != RetryNone && attempt+1 < policy.Attempts && policy.RetryableStatus != nil && policy.RetryableStatus(status)
}

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

func randomJitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	var seed [8]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return max / 2
	}
	return time.Duration(binary.LittleEndian.Uint64(seed[:]) % uint64(max+1))
}

func SanitizeErrorText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = secretQueryPattern.ReplaceAllString(text, "$1=[redacted]")
	text = bearerPattern.ReplaceAllString(text, "$1 [redacted]")
	const maxErrorTextLength = 512
	if len(text) <= maxErrorTextLength {
		return text
	}
	return strings.TrimSpace(text[:maxErrorTextLength]) + "..."
}

var (
	secretQueryPattern = regexp.MustCompile(`(?i)\b(api[_-]?key|token|password|secret|authorization)=([^&\s]+)`)
	bearerPattern      = regexp.MustCompile(`(?i)\b(Bearer|ApiKey)\s+[A-Za-z0-9._~+/=-]+`)
)
