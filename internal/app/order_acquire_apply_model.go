package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func providerOrderApplied(order core.Order, route core.Route, provider core.Provider, providerOrder core.ProviderOrder, now time.Time) core.Order {
	policy := providerPolicyForUpstreamOrder(provider, providerOrder.UpstreamOrderID).WithDefaults()
	acquiredAt := normalizedProviderAcquiredAt(providerOrder.AcquiredAt, now)
	expiresAt := normalizedProviderExpiresAt(order.ExpiresAt, providerOrder.ExpiresAt, acquiredAt, policy)
	order.ProviderKey = provider.Key()
	order.UpstreamOrderID = providerOrder.UpstreamOrderID
	order.Target = withRouteTargetDefaults(order.Target, route)
	order.PhoneNumber = providerOrder.PhoneNumber
	order.Status = core.StatusPendingCode
	order.Price = providerOrder.Price
	order.AcquiredAt = acquiredAt
	order.ExpiresAt = expiresAt
	order.UpdatedAt = now
	order.CancelAllowedAt = providerCancelAllowedAt(acquiredAt, policy)
	order.CanRequestAdditionalCode = providerOrder.CanRequestAdditionalCode
	order.LastError = nil
	return order
}
