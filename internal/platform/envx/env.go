package envx

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func String(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func StringDefault(name string, fallback string) string {
	if value := String(name); value != "" {
		return value
	}
	return fallback
}

func IntStrict(name string, fallback int) (int, error) {
	value := String(name)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	return parsed, nil
}
