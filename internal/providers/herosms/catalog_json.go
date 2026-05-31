package herosms

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func decodeHeroSMSJSONObject(result string, out any) error {
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return handlerapi.MapTextError(result)
	}
	return nil
}

func isHeroSMSUnsupportedCatalogLookup(err error) bool {
	providerErr, ok := err.(*core.Error)
	return ok && (providerErr.Code == core.CodeUnsupportedOperation || providerErr.Code == core.CodeValidationFailed || providerErr.Code == core.CodeUpstreamRejected)
}

func firstHeroSMSScalar(values ...json.RawMessage) string {
	for _, value := range values {
		if scalar := heroSMSScalar(value); scalar != "" {
			return scalar
		}
	}
	return ""
}

func heroSMSScalar(value json.RawMessage) string {
	if len(value) == 0 || string(value) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(value, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(value, &number); err == nil {
		return strings.TrimSpace(number.String())
	}
	return strings.Trim(string(value), "\"")
}

func heroSMSInt(value json.RawMessage) int {
	scalar := heroSMSScalar(value)
	if scalar == "" {
		return 0
	}
	integer, _ := strconv.Atoi(scalar)
	return integer
}
