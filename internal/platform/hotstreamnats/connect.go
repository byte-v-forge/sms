package hotstreamnats

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/nats-io/nats.go"
)

func ConnectService(ctx context.Context, cfg ServiceConfig, opts ...nats.Option) (*Bus, error) {
	if strings.TrimSpace(cfg.URL) == "" {
		message := strings.TrimSpace(cfg.RequiredMessage)
		if message == "" {
			message = "hotstream nats url is required"
		}
		return nil, errors.New(message)
	}
	clientName := strings.TrimSpace(cfg.ClientName)
	service := strings.TrimSpace(cfg.Service)
	if clientName == "" {
		clientName = service
	}
	if clientName == "" {
		clientName = "byte-v-forge"
	}
	subject := strings.TrimSpace(cfg.Subject)
	if subject == "" {
		subjectService := service
		if subjectService == "" {
			subjectService = clientName
		}
		subject = hotstream.ServiceStateSubject(subjectService)
	}
	return Connect(ctx, Config{
		URL:        cfg.URL,
		ClientName: clientName,
		Subject:    subject,
		BufferSize: cfg.BufferSize,
	}, opts...)
}

func Connect(ctx context.Context, cfg Config, opts ...nats.Option) (*Bus, error) {
	url := strings.TrimSpace(cfg.URL)
	if url == "" {
		return nil, errors.New("hotstream nats url is required")
	}
	name := strings.TrimSpace(cfg.ClientName)
	if name == "" {
		name = "byte-v-forge-hotstream"
	}
	subject := strings.TrimSpace(cfg.Subject)
	if subject == "" {
		subject = hotstream.ServiceStateSubject(name)
	}
	nodeID := nats.NewInbox()
	options := append([]nats.Option{
		nats.Name(name + " hotstream"),
		nats.Timeout(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
	}, opts...)
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, fmt.Errorf("connect hotstream nats: %w", err)
	}
	bus := &Bus{conn: conn, hub: hotstream.NewHub(cfg.BufferSize), subject: subject, nodeID: nodeID}
	sub, err := conn.Subscribe(subject, bus.receive)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("subscribe hotstream nats subject %s: %w", subject, err)
	}
	bus.sub = sub
	if err := conn.Flush(); err != nil {
		bus.Close()
		return nil, fmt.Errorf("flush hotstream nats subscription: %w", err)
	}
	go func() {
		<-ctx.Done()
		bus.Close()
	}()
	return bus, nil
}
