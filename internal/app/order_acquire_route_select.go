package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func acquireRequestRoute(order core.Order, route core.Route) core.Route {
	if strings.TrimSpace(route.ProviderKey) == "" {
		return routeFromOrder(order)
	}
	return route
}
