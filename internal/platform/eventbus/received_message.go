package eventbus

import (
	"context"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

type ReceivedMessage struct {
	Subject    string
	Envelope   *commonv1.EventEnvelope
	Extensions map[string]string
	Attempt    int32
	Ack        func(context.Context) error
	Nak        func(context.Context) error
	NakDelay   func(context.Context, time.Duration) error
	Term       func(context.Context) error
	DeadLetter func(context.Context, string) error
}
