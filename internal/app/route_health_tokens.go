package app

import (
	"strings"
	"time"
)

func normalizeRouteHealthToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func seconds(value time.Duration) int {
	if value <= 0 {
		return 0
	}
	return int(value.Round(time.Second) / time.Second)
}
