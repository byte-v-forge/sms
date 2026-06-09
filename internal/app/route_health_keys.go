package app

import (
	"strconv"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func routeHealthKey(route core.Route) string {
	policy := routeFailurePolicyWithDefaults(route.FailurePolicy)
	providerKey := normalizeProviderKey(route.ProviderKey)
	if providerKey == "" {
		return ""
	}
	parts := []string{
		normalizeRouteHealthToken(policy.ScopeKey),
		strconv.Itoa(policy.FailureThreshold),
		strconv.Itoa(seconds(policy.FailureWindow)),
		strconv.Itoa(seconds(policy.DisableTTL)),
		providerKey,
		normalizeRouteHealthToken(firstNonEmpty(route.UpstreamServiceKey, route.ApplicationKey)),
		routeCountryISO2(route.CountryISO2),
		routeCallingCode(route.CountryCallingCode),
		normalizeRouteHealthToken(route.ProviderCountryID),
		normalizeRouteHealthToken(route.UpstreamProviderID),
	}
	return strings.Join(parts, "\x1f")
}
