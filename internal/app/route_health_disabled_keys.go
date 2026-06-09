package app

import "github.com/byte-v-forge/sms/internal/core"

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
