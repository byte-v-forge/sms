package hotstream

import (
	"context"
	"errors"
	"sync"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

const DefaultBufferSize = 256

var ErrSlowConsumer = errors.New("hotstream slow consumer")

type Publisher interface {
	Publish(context.Context, *observabilityv1.HotStreamEvent) error
}

type Subscriber interface {
	Subscribe(context.Context, Filter) (*Subscription, error)
}

type Bus interface {
	Publisher
	Subscriber
}

type Hub struct {
	mu     sync.Mutex
	subs   map[*subscription]struct{}
	buffer int
}

func NewHub(buffer int) *Hub {
	if buffer <= 0 {
		buffer = DefaultBufferSize
	}
	return &Hub{subs: map[*subscription]struct{}{}, buffer: buffer}
}
