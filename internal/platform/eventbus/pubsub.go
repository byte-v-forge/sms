package eventbus

import "context"

type PublishAck struct {
	Stream    string
	Sequence  uint64
	Duplicate bool
}

type Publisher interface {
	Publish(context.Context, Message) (PublishAck, error)
}

type Consumer interface {
	Fetch(context.Context, int) ([]ReceivedMessage, error)
}
