package eventbus

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hashx"
)

func StableEventID(prefix string, parts ...string) string {
	return strings.TrimSpace(prefix) + hashx.StableParts(parts...)
}
