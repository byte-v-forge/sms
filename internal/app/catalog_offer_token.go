package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/searchx"
)

func catalogTokenMatches(candidate string, query string) bool {
	normalized := normalizeCatalogToken(candidate)
	return normalized != "" && (normalized == query || strings.Contains(normalized, query))
}

func normalizeCatalogToken(value string) string {
	return searchx.Token(value)
}
