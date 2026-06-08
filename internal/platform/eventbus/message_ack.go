package eventbus

import (
	"context"
	"time"
)

func Ack(ctx context.Context, action func(context.Context) error, label string, logf LogFunc) {
	if action == nil {
		return
	}
	if err := action(ctx); err != nil && ctx.Err() == nil {
		logger(logf)("%s failed: %v", label, err)
	}
}

func AckMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	Ack(ctx, message.Ack, label, logf)
}

func NakMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	Ack(ctx, message.Nak, label, logf)
}

func NakMessageDelay(ctx context.Context, message ReceivedMessage, delay time.Duration, label string, logf LogFunc) {
	if delay > 0 && message.NakDelay != nil {
		Ack(ctx, func(nakCtx context.Context) error { return message.NakDelay(nakCtx, delay) }, label, logf)
		return
	}
	NakMessage(ctx, message, label, logf)
}

func TermMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	if message.DeadLetter != nil {
		Ack(ctx, func(deadLetterCtx context.Context) error {
			return message.DeadLetter(deadLetterCtx, label)
		}, "publish dead letter for "+label, logf)
	}
	Ack(ctx, message.Term, label, logf)
}
