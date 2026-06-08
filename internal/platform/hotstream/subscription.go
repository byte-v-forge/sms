package hotstream

import (
	"context"
	"sync"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

type Subscription struct {
	Events <-chan *observabilityv1.HotStreamEvent
	hub    *Hub
	inner  *subscription
}

type subscription struct {
	filter Filter
	events chan *observabilityv1.HotStreamEvent
	done   chan struct{}
	err    error
	once   sync.Once
}

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

func (h *Hub) unsubscribe(sub *subscription, err error) {
	if h == nil || sub == nil {
		return
	}
	h.mu.Lock()
	delete(h.subs, sub)
	h.mu.Unlock()
	sub.close(err)
}

func (s *Subscription) Close() {
	if s == nil || s.inner == nil {
		return
	}
	if s.hub != nil {
		s.hub.unsubscribe(s.inner, nil)
		return
	}
	s.inner.close(nil)
}

func (s *Subscription) Err() error {
	if s == nil || s.inner == nil {
		return nil
	}
	return s.inner.err
}

func (s *Subscription) Done() <-chan struct{} {
	if s == nil || s.inner == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}
	return s.inner.done
}

func (s *subscription) close(err error) {
	s.once.Do(func() {
		s.err = err
		close(s.done)
		close(s.events)
	})
}
