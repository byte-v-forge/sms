package herosms

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func heroSMSBool(raw json.RawMessage) bool {
	scalar := strings.ToLower(jsonx.FirstScalar(raw))
	return scalar == "true" || scalar == "1"
}

func heroSMSCurrencyCode(raw json.RawMessage) string {
	switch jsonx.FirstScalar(raw) {
	case "840":
		return "USD"
	default:
		return ""
	}
}

func parseHeroSMSTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Time{}
}
