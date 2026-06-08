package smsbower

import (
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func providerTruthy(raw json.RawMessage) bool {
	switch strings.ToLower(jsonx.Scalar(raw)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}
