package natseventbus

import (
	"time"

	"github.com/nats-io/nats.go"
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
