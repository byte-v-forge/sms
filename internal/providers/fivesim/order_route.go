package fivesim

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func operatorForRoute(route core.Route) (string, error) {
	operator := strings.TrimSpace(route.UpstreamProviderID)
	if operator == "" {
		return "", core.NewError(core.CodeValidationFailed, "5sim operator is required", false)
	}
	return operator, nil
}
