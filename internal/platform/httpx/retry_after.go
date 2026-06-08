package httpx

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RetryAfter(header http.Header) time.Duration {
	value := strings.TrimSpace(header.Get("Retry-After"))
	if value == "" {
		return 0
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(value); err == nil {
		return time.Until(when)
	}
	return 0
}

func RetryAfterMax(header http.Header, maximum time.Duration) time.Duration {
	delay := RetryAfter(header)
	if delay <= 0 {
		return 0
	}
	if maximum > 0 && delay > maximum {
		return maximum
	}
	return delay
}
