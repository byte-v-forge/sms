package natseventbus

import (
	"context"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

func publishDeadLetter(ctx context.Context, bus *Bus, durable string, envelope *commonv1.EventEnvelope, attempt int32, reason string) error {
	if bus == nil || envelope == nil {
		return nil
	}
	message, err := deadLetterMessage(durable, envelope, attempt, reason)
	if err != nil {
		return err
	}
	_, err = bus.Publish(ctx, message)
	return err
}
