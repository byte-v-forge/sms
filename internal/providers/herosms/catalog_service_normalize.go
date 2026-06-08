package herosms

import (
	"strconv"
	"strings"
)

func normalizeHeroSMSServiceKey(value string) string {
	value = strings.TrimSpace(value)
	if index := strings.LastIndex(value, "_"); index > 0 {
		suffix := value[index+1:]
		if _, err := strconv.Atoi(suffix); err == nil {
			return value[:index]
		}
	}
	return value
}

func normalizeHeroSMSCatalogToken(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}
