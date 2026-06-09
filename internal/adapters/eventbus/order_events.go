package eventbusadapter

import (
	"context"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
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

func (b *OrderEventRecorder) AsyncRequests() bool { return b != nil }

func (b *OrderEventRecorder) record(_ context.Context, definition eventcatalog.Definition, message proto.Message, metadata *commonv1.EventMetadata, attrs map[string]string) (eventoutbox.Record, error) {
	if b == nil {
		return eventoutbox.Record{}, nil
	}
	return eventoutbox.NewRecordFor(definition, message, metadata, attrs)
}

func (b *OrderEventRecorder) metadata(eventName string, subject string, eventID string, correlationID string, occurredAt time.Time) *commonv1.EventMetadata {
	return newOrderEventMetadata(b.source, eventName, subject, eventID, correlationID, occurredAt)
}
