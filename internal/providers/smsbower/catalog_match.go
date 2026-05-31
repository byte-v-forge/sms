package smsbower

import "strings"

func matchService(candidate string, applications []ApplicationOffer) string {
	normalized := normalizeApplicationAlias(candidate)
	for _, app := range applications {
		if normalizeApplicationAlias(app.UpstreamServiceKey) == normalized || normalizeApplicationAlias(app.ApplicationKey) == normalized {
			return app.UpstreamServiceKey
		}
	}
	for _, app := range applications {
		display := normalizeApplicationAlias(app.DisplayName)
		if display != "" && (display == normalized || strings.Contains(display, normalized)) {
			return app.UpstreamServiceKey
		}
	}
	return ""
}

func normalizeApplicationAlias(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}
