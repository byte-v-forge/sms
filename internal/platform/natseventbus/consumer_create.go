package natseventbus

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
)

func (b *Bus) PullConsumer(cfg ConsumerConfig) (*PullConsumer, error) {
	if b == nil || b.js == nil {
		return nil, errors.New("nats event bus is not connected")
	}
	subject := strings.TrimSpace(cfg.Subject)
	if subject == "" {
		subject = DefaultSubject
	}
	durable := strings.TrimSpace(cfg.Durable)
	if durable == "" {
		return nil, errors.New("durable consumer name is required")
	}
	opts := []nats.SubOpt{
		nats.BindStream(normalizedStreamName(cfg.Stream)),
		nats.ManualAck(),
		nats.AckExplicit(),
	}
	if cfg.AckWait > 0 {
		opts = append(opts, nats.AckWait(cfg.AckWait))
	}
	sub, err := b.js.PullSubscribe(subject, durable, opts...)
	if err != nil {
		return nil, fmt.Errorf("create nats pull consumer %s: %w", durable, err)
	}
	batch := cfg.Batch
	if batch <= 0 {
		batch = 10
	}
	maxWait := cfg.MaxWait
	if maxWait <= 0 {
		maxWait = DefaultFetchWait
	}
	return &PullConsumer{sub: sub, bus: b, durable: durable, batch: batch, maxWait: maxWait}, nil
}
