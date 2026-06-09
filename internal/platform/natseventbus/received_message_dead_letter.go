package natseventbus

import (
	"context"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

func deadLetterPublisher(bus *Bus, durable string, envelope *commonv1.EventEnvelope, attempt int32) func(context.Context, string) error {
	return func(ctx context.Context, reason string) error {
		return publishDeadLetter(ctx, bus, durable, envelope, attempt, reason)
	}
}
