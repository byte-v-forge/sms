package eventbus

import (
	"context"
	"strings"

	"google.golang.org/protobuf/proto"
)

func handleTypedMessage[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T], received ReceivedMessage) {
	if err := cfg.Expected.ValidateReceived(received); err != nil {
		cfg.Logf("validate %s envelope failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return
	}
	message := cfg.NewMessage()
	if err := UnmarshalPayload(received, message); err != nil {
		cfg.Logf("decode %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return
	}
	if cfg.Validate != nil {
		if err := cfg.Validate(message); err != nil {
			cfg.Logf("validate %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
			TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
			return
		}
	}
	applyHandlerResult(ctx, received, cfg.Handler(ctx, message, received), cfg.Logf)
}

func applyHandlerResult(ctx context.Context, message ReceivedMessage, result HandlerResult, logf LogFunc) {
	label := strings.TrimSpace(result.Label)
	switch result.Action {
	case MessageActionNak:
		if label == "" {
			label = "nak event"
		}
		NakMessageDelay(ctx, message, result.Delay, label, logf)
	case MessageActionTerm:
		if label == "" {
			label = "terminate event"
		}
		TermMessage(ctx, message, label, logf)
	default:
		if label == "" {
			label = "ack event"
		}
		AckMessage(ctx, message, label, logf)
	}
}
