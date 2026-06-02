package eventbusadapter

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	commonv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

const defaultSourceService = "sms-service"

type OrderEventRecorder struct {
	source string
}

func NewOrderEventRecorder(source string) *OrderEventRecorder {
	source = strings.TrimSpace(source)
	if source == "" {
		source = defaultSourceService
	}
	return &OrderEventRecorder{source: source}
}

func (b *OrderEventRecorder) record(_ context.Context, definition eventcatalog.Definition, message proto.Message, eventCtx *commonv1.EventContext, attrs map[string]string) (eventoutbox.Record, error) {
	if b == nil {
		return eventoutbox.Record{}, nil
	}
	return eventoutbox.NewRecordFor(definition, message, eventCtx, attrs)
}

func (b *OrderEventRecorder) context(eventName string, eventID string, correlationID string, occurredAt time.Time) *commonv1.EventContext {
	return eventbus.NewEventContext(eventbus.EventContextConfig{
		EventID:       eventID,
		EventName:     eventName,
		EventVersion:  eventcatalog.EventVersionV1,
		OccurredAt:    occurredAt,
		SourceService: b.source,
		CorrelationID: correlationID,
	})
}

func orderAttributes(order core.Order) map[string]string {
	return eventbus.Attributes(
		"order_id", order.ID,
		"provider_key", order.ProviderKey,
		"status", string(order.Status),
	)
}
