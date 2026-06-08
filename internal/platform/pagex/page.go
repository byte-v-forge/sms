package pagex

import (
	"fmt"
	"strconv"
)

const (
	DefaultLimit = 100
	MaxLimit     = 500
)

func ParseOffsetToken(pageToken string) (int, error) {
	if pageToken == "" {
		return 0, nil
	}
	offset, err := strconv.Atoi(pageToken)
	if err != nil || offset < 0 {
		return 0, fmt.Errorf("page_token must be a non-negative offset")
	}
	return offset, nil
}

func OffsetToken(offset int) string {
	if offset <= 0 {
		return ""
	}
	return strconv.Itoa(offset)
}

func ClampSize(pageSize int, fallback int, maximum int) int {
	if fallback <= 0 {
		fallback = 50
	}
	if maximum <= 0 {
		maximum = fallback
	}
	if pageSize <= 0 {
		return fallback
	}
	if pageSize > maximum {
		return maximum
	}
	return pageSize
}

func NormalizePageLimit(limit int) int {
	return NormalizeLimit(limit, DefaultLimit, MaxLimit)
}

func NormalizeLimit(limit int, defaultLimit int, maxLimit int) int {
	if defaultLimit <= 0 {
		defaultLimit = DefaultLimit
	}
	if maxLimit <= 0 {
		maxLimit = defaultLimit
	}
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func TrimLimit[T any](rows []T, limit int) ([]T, bool) {
	if limit < 0 {
		limit = 0
	}
	if len(rows) <= limit {
		return rows, false
	}
	return rows[:limit], true
}
