package hotstream

import (
	"sync"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

type subscription struct {
	filter Filter
	events chan *observabilityv1.HotStreamEvent
	done   chan struct{}
	err    error
	once   sync.Once
}
