package hotstream

import "strings"

func cleanEventText(value string) string {
	return strings.TrimSpace(value)
}
