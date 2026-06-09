package natseventbus

import "github.com/nats-io/nats.go"

func deliveryAttempt(msg *nats.Msg) int32 {
	if msg == nil {
		return 0
	}
	meta, err := msg.Metadata()
	if err != nil || meta == nil {
		return 0
	}
	return int32(meta.NumDelivered)
}
