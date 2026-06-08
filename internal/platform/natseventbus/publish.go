package natseventbus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func (b *Bus) Publish(ctx context.Context, message eventbus.Message) (eventbus.PublishAck, error) {
	if b == nil || b.js == nil {
		return eventbus.PublishAck{}, errors.New("nats event bus is not connected")
	}
	envelope, err := eventbus.NewEnvelope(message)
	if err != nil {
		return eventbus.PublishAck{}, err
	}
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return eventbus.PublishAck{}, fmt.Errorf("marshal event envelope: %w", err)
	}
	msg := &nats.Msg{
		Subject: envelope.GetSubject(),
		Header:  envelopeHeaders(envelope),
		Data:    payload,
	}
	opts := []nats.PubOpt{nats.Context(ctx)}
	if idempotencyKey := strings.TrimSpace(envelope.GetMetadata().GetIdempotencyKey()); idempotencyKey != "" {
		opts = append(opts, nats.MsgId(idempotencyKey))
	}
	ack, err := b.js.PublishMsg(msg, opts...)
	if err != nil {
		return eventbus.PublishAck{}, fmt.Errorf("publish nats event %s: %w", envelope.GetSubject(), err)
	}
	return eventbus.PublishAck{
		Stream:    ack.Stream,
		Sequence:  ack.Sequence,
		Duplicate: ack.Duplicate,
	}, nil
}
