package app

import (
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func routeFailurePolicyIsZero(policy core.RouteFailurePolicy) bool {
	return strings.TrimSpace(policy.ScopeKey) == "" &&
		policy.FailureThreshold == 0 &&
		policy.FailureWindow <= 0 &&
		policy.DisableTTL <= 0
}

func secondsDuration(seconds int32) time.Duration {
	if seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}
