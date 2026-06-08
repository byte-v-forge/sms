package hotstreamnats

import (
	"context"
	"errors"
	"fmt"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

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
