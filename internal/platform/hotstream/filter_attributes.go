package hotstream

import "strings"

func matchAttributes(expected map[string]string, actual map[string]string) bool {
	if len(expected) == 0 {
		return true
	}
	for key, value := range expected {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if actual[strings.TrimSpace(key)] != strings.TrimSpace(value) {
			return false
		}
	}
	return true
}
