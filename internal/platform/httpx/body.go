package httpx

import "io"

func ReadLimited(body io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		limit = DefaultMaxBodyBytes
	}
	return io.ReadAll(io.LimitReader(body, limit))
}
