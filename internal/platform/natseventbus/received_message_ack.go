package natseventbus

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

func ackMessage(msg *nats.Msg) func(context.Context) error {
	return func(context.Context) error { return msg.Ack() }
}

func nakMessage(msg *nats.Msg) func(context.Context) error {
	return func(context.Context) error { return msg.Nak() }
}

func nakMessageDelay(msg *nats.Msg) func(context.Context, time.Duration) error {
	return func(_ context.Context, delay time.Duration) error { return msg.NakWithDelay(delay) }
}

func termMessage(msg *nats.Msg) func(context.Context) error {
	return func(context.Context) error { return msg.Term() }
}
