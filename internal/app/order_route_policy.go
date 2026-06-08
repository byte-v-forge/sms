package app

import (
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func routeFailurePolicyFromRoutePolicy(policy *smsv1.SmsRoutePolicy) core.RouteFailurePolicy {
	if policy == nil {
		return core.RouteFailurePolicy{}
	}
	return routeFailurePolicyFromProto(policy.GetFailurePolicy())
}

func routeFailurePolicyFromProto(policy *smsv1.SmsRouteFailurePolicy) core.RouteFailurePolicy {
	if policy == nil {
		return core.RouteFailurePolicy{}
	}
	return core.RouteFailurePolicy{
		ScopeKey:         strings.TrimSpace(policy.GetScopeKey()),
		FailureThreshold: int(policy.GetFailureThreshold()),
		FailureWindow:    secondsDuration(policy.GetFailureWindowSeconds()),
		DisableTTL:       secondsDuration(policy.GetDisableTtlSeconds()),
	}
}

func protoRouteFailurePolicy(policy core.RouteFailurePolicy) *smsv1.SmsRouteFailurePolicy {
	if routeFailurePolicyIsZero(policy) {
		return nil
	}
	return &smsv1.SmsRouteFailurePolicy{
		ScopeKey:             strings.TrimSpace(policy.ScopeKey),
		FailureThreshold:     int32(policy.FailureThreshold),
		FailureWindowSeconds: int32(durationSeconds(policy.FailureWindow)),
		DisableTtlSeconds:    int32(durationSeconds(policy.DisableTTL)),
	}
}

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

func durationSeconds(duration time.Duration) int {
	if duration <= 0 {
		return 0
	}
	return int(duration.Round(time.Second) / time.Second)
}
