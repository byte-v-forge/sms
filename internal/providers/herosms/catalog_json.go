package herosms

import (
	"encoding/json"

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
