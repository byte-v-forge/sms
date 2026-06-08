package hotstream

import (
	"context"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

func (h *Hub) Subscribe(ctx context.Context, filter Filter) (*Subscription, error) {
	if h == nil {
		h = NewHub(DefaultBufferSize)
	}
	sub := &subscription{
		filter: filter,
		events: make(chan *observabilityv1.HotStreamEvent, h.buffer),
		done:   make(chan struct{}),
	}
	h.mu.Lock()
	h.subs[sub] = struct{}{}
	h.mu.Unlock()
	go func() {
		<-ctx.Done()
		h.unsubscribe(sub, ctx.Err())
	}()
	return &Subscription{Events: sub.events, hub: h, inner: sub}, nil
}
