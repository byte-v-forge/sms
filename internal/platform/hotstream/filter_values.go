package hotstream

import "strings"

func matchAny(allowed []string, value string) bool {
	if len(allowed) == 0 {
		return true
	}
	value = strings.TrimSpace(value)
	for _, item := range allowed {
		if strings.TrimSpace(item) == value {
			return true
		}
	}
	return false
}
