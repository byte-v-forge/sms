package app

import (
	"strings"

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
		FailureWindowSeconds: int32(seconds(policy.FailureWindow)),
		DisableTtlSeconds:    int32(seconds(policy.DisableTTL)),
	}
}
