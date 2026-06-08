package hotstream

import observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"

type Subscription struct {
	Events <-chan *observabilityv1.HotStreamEvent
	hub    *Hub
	inner  *subscription
}
