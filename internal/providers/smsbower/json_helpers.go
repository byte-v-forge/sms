package smsbower

import (
	"encoding/json"
	"fmt"

	"github.com/byte-v-forge/sms/internal/core"
)

func decodeJSONObject(result string, out any) error {
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("bad json response: %v", err), false)
	}
	return nil
}
