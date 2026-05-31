package smsbower

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func rawJSONScalar(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String()
	}
	var floatValue float64
	if err := json.Unmarshal(raw, &floatValue); err == nil {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	return strings.Trim(string(raw), "\"")
}
func decodeJSONObject(result string, out any) error {
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("bad json response: %v", err), false)
	}
	return nil
}
