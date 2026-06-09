package eventbus

import (
	"context"

	"google.golang.org/protobuf/proto"
)

func handleTypedMessage[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T], received ReceivedMessage) {
	message, ok := decodeTypedMessage(ctx, cfg, received)
	if !ok {
		return
	}
	applyHandlerResult(ctx, received, cfg.Handler(ctx, message, received), cfg.Logf)
}
