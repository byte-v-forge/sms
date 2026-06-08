package natseventbus

import (
	"time"

	"github.com/nats-io/nats.go"
)

func natsOptions(clientName string, opts []nats.Option) []nats.Option {
	return append([]nats.Option{
		nats.Name(clientName),
		nats.Timeout(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
	}, opts...)
}
