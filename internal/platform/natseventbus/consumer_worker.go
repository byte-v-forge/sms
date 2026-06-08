package natseventbus

import (
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
)

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
