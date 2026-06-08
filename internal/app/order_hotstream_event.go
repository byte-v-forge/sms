package app

import (
	"fmt"
	"time"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func hotStreamOrderEvent(order core.Order) *observabilityv1.HotStreamEvent {
	updatedAt := order.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	return hotstream.NewEvent(hotstream.EventConfig{
		EventID:       eventbus.StableEventID("sms-order-", order.ID, string(order.Status), fmt.Sprintf("%d", updatedAt.UnixNano())),
		EventType:     SMSOrderUpdatedEvent,
		SourceService: SMSHotStreamSource,
		ResourceType:  SMSOrderResource,
		ResourceID:    order.ID,
		Scope:         string(order.Status),
		OccurredAt:    updatedAt,
		CorrelationID: order.RequestID,
		Attributes: map[string]string{
			"order_id":     order.ID,
			"request_id":   order.RequestID,
			"provider_key": order.ProviderKey,
			"status":       string(order.Status),
		},
	})
}
