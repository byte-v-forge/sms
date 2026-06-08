package secretref

import (
	"fmt"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hashx"
)

func StableID(prefix string, parts ...string) string {
	prefix = cleanSegment(prefix)
	if prefix == "" {
		prefix = "secret"
	}
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			clean = append(clean, value)
		}
	}
	if len(clean) == 0 {
		return prefix
	}
	return fmt.Sprintf("%s-%s", prefix, hashx.ShortSHA256(hashx.StableParts(clean...), 24))
}

func cleanSegment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(" ", "-", "_", "-", "/", "-", ":", "-").Replace(value)
	value = strings.Trim(value, "-")
	return value
}
