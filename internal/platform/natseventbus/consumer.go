package natseventbus

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ConsumerConfig struct {
	Stream  string
	Subject string
	Durable string
	Batch   int
	MaxWait time.Duration
	AckWait time.Duration
}

type PullConsumer struct {
	sub     *nats.Subscription
	bus     *Bus
	durable string
	batch   int
	maxWait time.Duration
}

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

func (b *Bus) PullWorkerConsumer(stream string, subject string, durable string, batch int, ackWait time.Duration) (*PullConsumer, error) {
	return b.PullConsumer(ConsumerConfig{
		Stream:  stream,
		Subject: subject,
		Durable: durable,
		Batch:   batch,
		MaxWait: DefaultFetchWait,
		AckWait: ackWait,
	})
}

func (b *Bus) PullWorkerForDefinition(stream string, definition eventcatalog.Definition, batch int, ackWait time.Duration) (*PullConsumer, error) {
	return b.PullWorkerForBinding(stream, definition.DefaultConsumerBinding(), batch, ackWait)
}

func (b *Bus) PullWorkerForBinding(stream string, binding eventcatalog.ConsumerBinding, batch int, ackWait time.Duration) (*PullConsumer, error) {
	if err := binding.Validate(); err != nil {
		return nil, err
	}
	return b.PullWorkerConsumer(stream, binding.Subject(), binding.DurableName(), batch, ackWait)
}

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
