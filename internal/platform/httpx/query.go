package httpx

import (
	"net/http"
	"strconv"
	"strings"
)

func QueryInt(r *http.Request, key string, fallback int) int {
	if r == nil {
		return fallback
	}
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func QueryBool(r *http.Request, key string, fallback bool) bool {
	if r == nil {
		return fallback
	}
	value := strings.ToLower(strings.TrimSpace(r.URL.Query().Get(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes"
}
