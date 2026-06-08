package hotstreamnats

import (
	"context"
	"fmt"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/nats-io/nats.go"
)

func ConnectService(ctx context.Context, cfg ServiceConfig, opts ...nats.Option) (*Bus, error) {
	connectConfig, err := normalizeServiceConfig(cfg)
	if err != nil {
		return nil, err
	}
	return Connect(ctx, connectConfig, opts...)
}

func Connect(ctx context.Context, cfg Config, opts ...nats.Option) (*Bus, error) {
	cfg, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}
	conn, err := nats.Connect(cfg.URL, hotStreamNATSOptions(cfg.ClientName, opts)...)
	if err != nil {
		return nil, fmt.Errorf("connect hotstream nats: %w", err)
	}
	bus := &Bus{conn: conn, hub: hotstream.NewHub(cfg.BufferSize), subject: cfg.Subject, nodeID: nats.NewInbox()}
	if err := bus.subscribeConn(); err != nil {
		conn.Close()
		return nil, err
	}
	closeBusOnContext(ctx, bus)
	return bus, nil
}
