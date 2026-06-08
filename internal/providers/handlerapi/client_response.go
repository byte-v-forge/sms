package handlerapi

import (
	"fmt"

	"github.com/byte-v-forge/sms/internal/core"
)

func handlerAPIHTTPError(statusCode int, text string) error {
	if text != "" {
		return MapTextError(text)
	}
	return core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("handler api http status %d", statusCode), true)
}
