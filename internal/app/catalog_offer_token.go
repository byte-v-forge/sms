package app

import "strings"

func catalogTokenMatches(candidate string, query string) bool {
	normalized := normalizeCatalogToken(candidate)
	return normalized != "" && (normalized == query || strings.Contains(normalized, query))
}

func normalizeCatalogToken(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}
