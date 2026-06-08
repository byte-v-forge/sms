package httpx

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultMaxBodyBytes int64 = 8 * 1024 * 1024

func ReadLimited(body io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		limit = DefaultMaxBodyBytes
	}
	return io.ReadAll(io.LimitReader(body, limit))
}

func ReadMaybeGzipLimited(body io.Reader, limit int64) ([]byte, error) {
	raw, err := ReadLimited(body, limit)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(raw, []byte{0x1f, 0x8b}) {
		return raw, nil
	}
	reader, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ReadLimited(reader, limit)
}

func QueryInt(r *http.Request, key string, fallback int) int {
	if r == nil {
		return fallback
	}
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func Successful(status int) bool {
	return status >= 200 && status < 300
}

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

func QueryBool(r *http.Request, key string, fallback bool) bool {
	if r == nil {
		return fallback
	}
	value := strings.ToLower(strings.TrimSpace(r.URL.Query().Get(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes"
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
