package natseventbus

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Bus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func Connect(cfg Config, opts ...nats.Option) (*Bus, error) {
	cfg = cfg.withDefaults()
	conn, err := nats.Connect(cfg.URL, natsOptions(cfg.ClientName, opts)...)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("initialize jetstream: %w", err)
	}
	return &Bus{conn: conn, js: js}, nil
}
