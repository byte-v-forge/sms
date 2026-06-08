package hotstreamnats

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

type Config struct {
	URL        string
	ClientName string
	Subject    string
	BufferSize int
}

type ServiceConfig struct {
	URL             string
	Service         string
	ClientName      string
	Subject         string
	BufferSize      int
	RequiredMessage string
}

type Bus struct {
	conn    *nats.Conn
	hub     *hotstream.Hub
	subject string
	nodeID  string
	sub     *nats.Subscription
}

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

func (b *Bus) Close() {
	if b == nil {
		return
	}
	if b.sub != nil {
		_ = b.sub.Unsubscribe()
	}
	if b.conn != nil {
		b.conn.Drain()
		b.conn.Close()
	}
}

func (b *Bus) Publish(ctx context.Context, event *observabilityv1.HotStreamEvent) error {
	if b == nil || b.hub == nil {
		return errors.New("hotstream bus is not configured")
	}
	if event == nil {
		return nil
	}
	if err := b.hub.Publish(ctx, event); err != nil {
		return err
	}
	if b.conn == nil {
		return nil
	}
	payload, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal hotstream event: %w", err)
	}
	msg := &nats.Msg{Subject: b.subject, Header: nats.Header{}, Data: payload}
	msg.Header.Set("Bvf-Hotstream-Node", b.nodeID)
	msg.Header.Set("Bvf-Hotstream-Event-Type", event.GetMetadata().GetType())
	msg.Header.Set("Bvf-Hotstream-Resource-Type", event.GetResourceType())
	msg.Header.Set("Bvf-Hotstream-Resource-Id", event.GetResourceId())
	return b.conn.PublishMsg(msg)
}

func (b *Bus) Subscribe(ctx context.Context, filter hotstream.Filter) (*hotstream.Subscription, error) {
	if b == nil || b.hub == nil {
		return nil, errors.New("hotstream bus is not configured")
	}
	return b.hub.Subscribe(ctx, filter)
}

func (b *Bus) receive(msg *nats.Msg) {
	if b == nil || msg == nil || msg.Header.Get("Bvf-Hotstream-Node") == b.nodeID {
		return
	}
	event := &observabilityv1.HotStreamEvent{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		return
	}
	_ = b.hub.Publish(context.Background(), event)
}
