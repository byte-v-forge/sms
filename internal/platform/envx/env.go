package envx

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
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

func Bool(name string, fallback bool) bool {
	return ParseBool(String(name), fallback)
}

func ParseBool(value string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func Int(name string, fallback int) int {
	value := String(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
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

func PositiveInt(name string, fallback int) int {
	parsed := Int(name, fallback)
	if parsed <= 0 {
		return fallback
	}
	return parsed
}

func NonNegativeInt(name string, fallback int) int {
	parsed := Int(name, fallback)
	if parsed < 0 {
		return fallback
	}
	return parsed
}

func PositiveInt32(name string, fallback int32) int32 {
	parsed := PositiveInt(name, int(fallback))
	return int32(parsed)
}

func DurationSeconds(name string, fallback time.Duration) time.Duration {
	value := String(name)
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

func PositiveDurationSeconds(name string, fallback time.Duration) time.Duration {
	seconds := PositiveInt(name, int(fallback/time.Second))
	return time.Duration(seconds) * time.Second
}

func NonNegativeDurationSeconds(name string, fallback time.Duration) time.Duration {
	seconds := NonNegativeInt(name, int(fallback/time.Second))
	return time.Duration(seconds) * time.Second
}

func FloatDurationSeconds(name string, fallback float64) time.Duration {
	value := String(name)
	if value == "" {
		return time.Duration(fallback * float64(time.Second))
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed < 0 {
		return time.Duration(fallback * float64(time.Second))
	}
	return time.Duration(parsed * float64(time.Second))
}

func List(name string) []string {
	value := String(name)
	if value == "" {
		return nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			items = append(items, item)
		}
	}
	return items
}

func JSONStringMap(name string) (map[string]string, error) {
	value := String(name)
	if value == "" {
		return nil, nil
	}
	items := map[string]string{}
	if err := json.Unmarshal([]byte(value), &items); err != nil {
		return nil, fmt.Errorf("%s must be a JSON object with string values: %w", name, err)
	}
	normalized := make(map[string]string, len(items))
	for key, item := range items {
		key = strings.TrimSpace(key)
		item = strings.TrimSpace(item)
		if key == "" {
			return nil, fmt.Errorf("%s contains an empty key", name)
		}
		if item == "" {
			return nil, fmt.Errorf("%s contains an empty value for key %q", name, key)
		}
		normalized[key] = item
	}
	return normalized, nil
}
