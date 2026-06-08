package natseventbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/nats-io/nats.go"
)

func (c *PullConsumer) Fetch(ctx context.Context, batch int) ([]eventbus.ReceivedMessage, error) {
	if c == nil || c.sub == nil {
		return nil, errors.New("nats pull consumer is not configured")
	}
	if batch <= 0 {
		batch = c.batch
	}
	fetchCtx, cancel := context.WithTimeout(ctx, c.maxWait)
	defer cancel()
	messages, err := c.sub.Fetch(batch, nats.Context(fetchCtx))
	if errors.Is(err, nats.ErrTimeout) || errors.Is(err, context.DeadlineExceeded) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fetch nats messages: %w", err)
	}
	out := make([]eventbus.ReceivedMessage, 0, len(messages))
	for _, msg := range messages {
		received, err := receivedMessage(c.bus, c.durable, msg)
		if err != nil {
			_ = msg.Term()
			return nil, err
		}
		out = append(out, received)
	}
	return out, nil
}
