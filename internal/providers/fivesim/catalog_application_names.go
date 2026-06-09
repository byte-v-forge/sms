package fivesim

import (
	"strings"
	"unicode"
)

func fiveSimApplicationName(product string) string {
	return titleApplicationKey(product)
}

func titleApplicationKey(value string) string {
	parts := strings.FieldsFunc(strings.TrimSpace(value), func(r rune) bool {
		return r == '-' || r == '_' || r == '.' || unicode.IsSpace(r)
	})
	for index, part := range parts {
		parts[index] = titleApplicationPart(part)
	}
	return strings.Join(parts, " ")
}

func titleApplicationPart(value string) string {
	if value == "" || hasMixedCase(value) {
		return value
	}
	runes := []rune(strings.ToLower(value))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func hasMixedCase(value string) bool {
	return value != strings.ToLower(value) && value != strings.ToUpper(value)
}
