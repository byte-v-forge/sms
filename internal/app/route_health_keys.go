package app

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func disabledRouteRedisKeys(routes []core.Route) ([]string, []string) {
	seen := map[string]struct{}{}
	routeKeys := make([]string, 0, len(routes))
	redisKeys := make([]string, 0, len(routes))
	for _, route := range routes {
		routeKey := routeHealthKey(route)
		if routeKey == "" {
			continue
		}
		if _, exists := seen[routeKey]; exists {
			continue
		}
		seen[routeKey] = struct{}{}
		routeKeys = append(routeKeys, routeKey)
		redisKeys = append(redisKeys, routeHealthRedisKey("disabled", routeKey))
	}
	return routeKeys, redisKeys
}

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

func normalizeRouteHealthToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func routeHealthRedisKey(kind string, routeKey string) string {
	digest := sha256.Sum256([]byte(routeKey))
	return "sms:route_health:" + kind + ":" + hex.EncodeToString(digest[:])
}

func seconds(value time.Duration) int {
	if value <= 0 {
		return 0
	}
	return int(value.Round(time.Second) / time.Second)
}
