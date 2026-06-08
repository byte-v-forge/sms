package natseventbus

import (
	"fmt"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/nats-io/nats.go"
)

const (
	DefaultURL       = nats.DefaultURL
	DefaultStream    = eventcatalog.StreamName
	DefaultSubject   = eventcatalog.StreamSubject
	DefaultFetchWait = 5 * time.Second
)

type Config struct {
	URL        string
	ClientName string
}

type Bus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func Connect(cfg Config, opts ...nats.Option) (*Bus, error) {
	url := strings.TrimSpace(cfg.URL)
	if url == "" {
		url = DefaultURL
	}
	name := strings.TrimSpace(cfg.ClientName)
	if name == "" {
		name = "byte-v-forge"
	}
	options := append([]nats.Option{
		nats.Name(name),
		nats.Timeout(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
	}, opts...)
	conn, err := nats.Connect(url, options...)
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

func (b *Bus) Close() {
	if b == nil || b.conn == nil {
		return
	}
	b.conn.Drain()
	b.conn.Close()
}

func normalizedStreamName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultStream
	}
	return value
}
