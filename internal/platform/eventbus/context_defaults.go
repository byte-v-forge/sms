package eventbus

import (
	"strings"
	"time"
)

func eventMetadataVersion(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultEventVersion
	}
	return value
}

func eventMetadataTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now()
	}
	return value
}

func eventMetadataIdempotencyKey(cfg EventMetadataConfig) string {
	idempotencyKey := strings.TrimSpace(cfg.IdempotencyKey)
	if idempotencyKey == "" {
		return strings.TrimSpace(cfg.EventID)
	}
	return idempotencyKey
}
