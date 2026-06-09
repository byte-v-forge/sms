package geox

import (
	"regexp"
	"strings"
)

var countryNameSeparator = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeCountryNameText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = countryNameSeparator.ReplaceAllString(value, " ")
	return strings.Join(strings.Fields(value), " ")
}
