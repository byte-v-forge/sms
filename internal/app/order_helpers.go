package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) provider(key string) (core.Provider, error) {
	provider, ok := s.providers[key]
	if !ok {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider not registered", false)
	}
	return provider, nil
}

func (s *OrderService) routePolicy(ctx context.Context, route core.Route) core.ProviderPolicy {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.ProviderPolicy{}.WithDefaults()
	}
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.LoadPolicyForOrder(ctx, "").WithDefaults()
	}
	return provider.Policy().WithDefaults()
}

func withRouteTargetDefaults(target core.Target, route core.Route) core.Target {
	if target.ApplicationKey == "" {
		target.ApplicationKey = route.ApplicationKey
	}
	if target.CountryISO2 == "" {
		target.CountryISO2 = route.CountryISO2
	}
	if target.CountryCallingCode == "" {
		target.CountryCallingCode = route.CountryCallingCode
	}
	return target
}

func routeFromOrder(order core.Order) core.Route {
	return core.Route{
		ProviderKey:        order.ProviderKey,
		ApplicationKey:     order.Target.ApplicationKey,
		UpstreamServiceKey: order.Target.ApplicationKey,
		CountryISO2:        order.Target.CountryISO2,
		CountryCallingCode: order.Target.CountryCallingCode,
	}
}

func orderRequestExpiresAt(now time.Time, policy core.ProviderPolicy, lease time.Duration) time.Time {
	policy = policy.WithDefaults()
	ttl := policy.OrderTTL
	if lease > 0 && lease < ttl {
		ttl = lease
	}
	if ttl <= 0 {
		return time.Time{}
	}
	return now.Add(ttl)
}

func remainingLease(now time.Time, expiresAt time.Time) time.Duration {
	if expiresAt.IsZero() {
		return 0
	}
	if !expiresAt.After(now) {
		return 0
	}
	return expiresAt.Sub(now)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func asCoreError(err error) *core.Error {
	if err == nil {
		return nil
	}
	if smsErr, ok := err.(*core.Error); ok {
		return smsErr
	}
	if smsErr := runtimeCoreError(err); smsErr != nil {
		return smsErr
	}
	return core.NewError(core.CodeInternal, "sms service request failed", false)
}
