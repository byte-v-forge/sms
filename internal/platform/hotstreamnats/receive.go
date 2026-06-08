package hotstreamnats

import (
	"context"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

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
