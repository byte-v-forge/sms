package natseventbus

import (
	"context"
	"fmt"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
)

func publishDeadLetter(ctx context.Context, bus *Bus, durable string, envelope *commonv1.EventEnvelope, attempt int32, reason string) error {
	if bus == nil || envelope == nil {
		return nil
	}
	original := envelope.GetMetadata()
	originalID := ""
	originalName := ""
	originalVersion := ""
	originalSource := ""
	correlationID := ""
	traceID := ""
	if original != nil {
		originalID = original.GetId()
		originalName = original.GetType()
		originalVersion = original.GetVersion()
		originalSource = original.GetSource()
		correlationID = original.GetCorrelationId()
		traceID = original.GetTraceId()
	}
	eventID := eventbus.StableEventID("dead-letter-", envelope.GetSubject(), originalID, durable, fmt.Sprintf("%d", attempt))
	deadMetadata := eventbus.NewEventMetadata(eventbus.EventMetadataConfig{
		EventID:       eventID,
		EventName:     "platform.dead_letter",
		EventVersion:  eventcatalog.EventVersionV1,
		SourceService: "platform-eventbus",
		Subject:       eventcatalog.DeadLetter.Subject,
		CorrelationID: correlationID,
		TraceID:       traceID,
	})
	message, err := eventcatalog.DeadLetter.NewMessage(
		&commonv1.DeadLetterEvent{
			Metadata:             deadMetadata,
			OriginalSubject:      envelope.GetSubject(),
			OriginalEventId:      originalID,
			OriginalEventType:    originalName,
			OriginalEventVersion: originalVersion,
			OriginalSource:       originalSource,
			ConsumerDurable:      durable,
			DeliveryAttempt:      attempt,
			ErrorCode:            "terminated",
			ErrorMessage:         reason,
			CorrelationId:        correlationID,
		},
		deadMetadata,
		eventbus.Attributes(
			"original_subject", envelope.GetSubject(),
			"original_event_id", originalID,
			"consumer_durable", durable,
			"delivery_attempt", fmt.Sprintf("%d", attempt),
		),
	)
	if err != nil {
		return err
	}
	_, err = bus.Publish(ctx, message)
	return err
}
