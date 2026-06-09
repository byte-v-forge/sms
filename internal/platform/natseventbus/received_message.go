package natseventbus

import (
	"fmt"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func receivedMessage(bus *Bus, durable string, msg *nats.Msg) (eventbus.ReceivedMessage, error) {
	envelope := &commonv1.EventEnvelope{}
	if err := proto.Unmarshal(msg.Data, envelope); err != nil {
		return eventbus.ReceivedMessage{}, fmt.Errorf("decode nats event envelope: %w", err)
	}
	attempt := deliveryAttempt(msg)
	return eventbus.ReceivedMessage{
		Subject:    msg.Subject,
		Envelope:   envelope,
		Extensions: envelope.GetExtensions(),
		Attempt:    attempt,
		Ack:        ackMessage(msg),
		Nak:        nakMessage(msg),
		NakDelay:   nakMessageDelay(msg),
		Term:       termMessage(msg),
		DeadLetter: deadLetterPublisher(bus, durable, envelope, attempt),
	}, nil
}
