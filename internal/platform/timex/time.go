package timex

import (
	"strings"
	"time"
)

func ParseRFC3339(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func Unix(value string) int64 {
	parsed, ok := ParseRFC3339(value)
	if !ok {
		return 0
	}
	return parsed.Unix()
}

func UnixNano(value string) int64 {
	parsed, ok := ParseRFC3339(value)
	if !ok {
		return 0
	}
	return parsed.UnixNano()
}

func UnixFloat(value string) float64 {
	parsed, ok := ParseRFC3339(value)
	if !ok {
		return 0
	}
	return float64(parsed.UnixNano()) / float64(time.Second)
}
