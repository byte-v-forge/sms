package hotstream

import (
	"context"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
	"google.golang.org/protobuf/proto"
)

func (h *Hub) Publish(_ context.Context, event *observabilityv1.HotStreamEvent) error {
	if h == nil || event == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for sub := range h.subs {
		if !sub.filter.Match(event) {
			continue
		}
		cloned, _ := proto.Clone(event).(*observabilityv1.HotStreamEvent)
		if cloned == nil {
			continue
		}
		select {
		case sub.events <- cloned:
		default:
			sub.close(ErrSlowConsumer)
			delete(h.subs, sub)
		}
	}
	return nil
}
