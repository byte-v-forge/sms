package stringx

import (
	"fmt"
	"strings"
)

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func Digits(value string) string {
	var out strings.Builder
	for _, ch := range value {
		if ch >= '0' && ch <= '9' {
			out.WriteRune(ch)
		}
	}
	return out.String()
}

func ContainsFold(value string, keyword string) bool {
	if keyword == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(keyword))
}

func CompactWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func CompactSnippet(value string, limit int) string {
	value = CompactWhitespace(value)
	if value == "" {
		return ""
	}
	if limit > 0 && len(value) > limit {
		return value[:limit] + "..."
	}
	return value
}

func FirstNonEmptyAny(values ...any) string {
	for _, value := range values {
		if text := strings.TrimSpace(fmt.Sprint(value)); text != "" && text != "<nil>" {
			return text
		}
	}
	return ""
}
