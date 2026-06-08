package handlerapi

import "strings"

func normalizeHandlerAPIError(text string) (string, string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	if code, message, ok := parseHandlerAPIJSONError(text); ok {
		return code, message
	}
	code := text
	if idx := strings.Index(code, ":"); idx >= 0 {
		code = strings.TrimSpace(code[:idx])
	}
	return code, truncateHandlerAPIErrorMessage(text)
}
