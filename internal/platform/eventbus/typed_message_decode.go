package eventbus

import (
	"context"

	"google.golang.org/protobuf/proto"
)

func decodeTypedMessage[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T], received ReceivedMessage) (T, bool) {
	var zero T
	if err := cfg.Expected.ValidateReceived(received); err != nil {
		cfg.Logf("validate %s envelope failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return zero, false
	}
	message := cfg.NewMessage()
	if err := UnmarshalPayload(received, message); err != nil {
		cfg.Logf("decode %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return zero, false
	}
	if err := validateTypedMessage(cfg, message); err != nil {
		cfg.Logf("validate %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return zero, false
	}
	return message, true
}

func validateTypedMessage[T proto.Message](cfg TypedConsumerWorkerConfig[T], message T) error {
	if cfg.Validate == nil {
		return nil
	}
	return cfg.Validate(message)
}
