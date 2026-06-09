package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func validateAcquireRoute(route core.Route) error {
	if strings.TrimSpace(route.ProviderKey) == "" || strings.TrimSpace(route.UpstreamServiceKey) == "" || strings.TrimSpace(route.ProviderCountryID) == "" {
		return core.NewError(core.CodeValidationFailed, "sms acquire params are incomplete", false)
	}
	if providerRequiresUpstreamProviderID(route.ProviderKey) && strings.TrimSpace(route.UpstreamProviderID) == "" {
		return core.NewError(core.CodeValidationFailed, "sms upstream provider id is required", false)
	}
	return nil
}

func providerRequiresUpstreamProviderID(providerKey string) bool {
	switch normalizeProviderKey(providerKey) {
	case "5sim", "smsbower":
		return true
	default:
		return false
	}
}
