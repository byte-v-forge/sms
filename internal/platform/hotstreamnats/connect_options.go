package hotstreamnats

import (
	"time"

	"github.com/nats-io/nats.go"
)

func hotStreamNATSOptions(clientName string, opts []nats.Option) []nats.Option {
	options := append([]nats.Option{
		nats.Name(clientName + " hotstream"),
		nats.Timeout(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
	}, opts...)
	return options
}
