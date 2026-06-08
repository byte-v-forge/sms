package providerhttp

import (
	"regexp"
	"strings"
)

func SanitizeErrorText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = secretQueryPattern.ReplaceAllString(text, "$1=[redacted]")
	text = bearerPattern.ReplaceAllString(text, "$1 [redacted]")
	const maxErrorTextLength = 512
	if len(text) <= maxErrorTextLength {
		return text
	}
	return strings.TrimSpace(text[:maxErrorTextLength]) + "..."
}

var (
	secretQueryPattern = regexp.MustCompile(`(?i)\b(api[_-]?key|token|password|secret|authorization)=([^&\s]+)`)
	bearerPattern      = regexp.MustCompile(`(?i)\b(Bearer|ApiKey)\s+[A-Za-z0-9._~+/=-]+`)
)
