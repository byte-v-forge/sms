package app

import (
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

const (
	defaultRouteFailureThreshold = 3
	defaultRouteFailureWindow    = 10 * time.Minute
	defaultRouteDisableTTL       = 10 * time.Minute
)

func routeFailureCountsTowardDisable(err *core.Error) bool {
	if err == nil || !err.Retryable {
		return false
	}
	switch err.Code {
	case core.CodeNoNumberAvailable, core.CodeRateLimited, core.CodeSupplyUnavailable, core.CodeTimeout, core.CodeUpstreamRejected:
		return true
	default:
		return false
	}
}

func routeFailurePolicyWithDefaults(policy core.RouteFailurePolicy) core.RouteFailurePolicy {
	if policy.FailureThreshold <= 0 {
		policy.FailureThreshold = defaultRouteFailureThreshold
	}
	if policy.FailureWindow <= 0 {
		policy.FailureWindow = defaultRouteFailureWindow
	}
	if policy.DisableTTL <= 0 {
		policy.DisableTTL = defaultRouteDisableTTL
	}
	policy.ScopeKey = strings.TrimSpace(policy.ScopeKey)
	return policy
}
