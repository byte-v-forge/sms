package eventbusadapter

import (
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func orderAttributes(order core.Order) map[string]string {
	return eventbus.Attributes(
		"order_id", order.ID,
		"provider_key", order.ProviderKey,
		"status", string(order.Status),
	)
}
