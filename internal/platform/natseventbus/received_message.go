package natseventbus

import (
	"context"
	"fmt"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func receivedMessage(bus *Bus, durable string, msg *nats.Msg) (eventbus.ReceivedMessage, error) {
	envelope := &commonv1.EventEnvelope{}
	if err := proto.Unmarshal(msg.Data, envelope); err != nil {
		return eventbus.ReceivedMessage{}, fmt.Errorf("decode nats event envelope: %w", err)
	}
	attempt := deliveryAttempt(msg)
	return eventbus.ReceivedMessage{
		Subject:    msg.Subject,
		Envelope:   envelope,
		Extensions: envelope.GetExtensions(),
		Attempt:    attempt,
		Ack: func(context.Context) error {
			return msg.Ack()
		},
		Nak: func(context.Context) error {
			return msg.Nak()
		},
		NakDelay: func(_ context.Context, delay time.Duration) error {
			return msg.NakWithDelay(delay)
		},
		Term: func(context.Context) error {
			return msg.Term()
		},
		DeadLetter: func(ctx context.Context, reason string) error {
			return publishDeadLetter(ctx, bus, durable, envelope, attempt, reason)
		},
	}, nil
}

func deliveryAttempt(msg *nats.Msg) int32 {
	if msg == nil {
		return 0
	}
	meta, err := msg.Metadata()
	if err != nil || meta == nil {
		return 0
	}
	return int32(meta.NumDelivered)
}
