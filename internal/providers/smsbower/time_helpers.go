package smsbower

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func parseOrderTime(raw json.RawMessage) time.Time {
	return parseOrderTimeText(jsonx.Scalar(raw))
}

func parseOrderTimeText(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC()
	}
	if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
		if unix > 1_000_000_000_000 {
			return time.UnixMilli(unix).UTC()
		}
		return time.Unix(unix, 0).UTC()
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Now().UTC()
}
