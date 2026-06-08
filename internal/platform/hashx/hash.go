package hashx

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func ShortSHA256(value string, length int) string {
	text := SHA256Hex(value)
	if length <= 0 || length >= len(text) {
		return text
	}
	return text[:length]
}

func StableParts(parts ...string) string {
	return SHA256Hex(strings.Join(parts, "\x00"))
}
