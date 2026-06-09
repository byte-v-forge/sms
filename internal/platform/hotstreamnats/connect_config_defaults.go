package hotstreamnats

import "strings"

func requiredServiceConfigMessage(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "hotstream nats url is required"
	}
	return value
}
