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

func (b *OrderEventRecorder) record(_ context.Context, definition eventcatalog.Definition, message proto.Message, metadata *commonv1.EventMetadata, attrs map[string]string) (eventoutbox.Record, error) {
	if b == nil {
		return eventoutbox.Record{}, nil
	}
	return eventoutbox.NewRecordFor(definition, message, metadata, attrs)
}

func (b *OrderEventRecorder) metadata(eventName string, subject string, eventID string, correlationID string, occurredAt time.Time) *commonv1.EventMetadata {
	return eventbus.NewEventMetadata(eventbus.EventMetadataConfig{
		EventID:       eventID,
		EventName:     eventName,
		EventVersion:  eventcatalog.EventVersionV1,
		OccurredAt:    occurredAt,
		SourceService: b.source,
		Subject:       subject,
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
