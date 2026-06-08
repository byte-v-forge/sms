package eventbus

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
)

type TypedMessageHandler[T proto.Message] func(context.Context, T, ReceivedMessage) HandlerResult

type TypedConsumerWorkerConfig[T proto.Message] struct {
	Name            string
	Consumer        Consumer
	Expected        ExpectedMessage
	NewMessage      func() T
	Validate        func(T) error
	Handler         TypedMessageHandler[T]
	MalformedLabel  string
	Batch           int
	FetchErrorDelay time.Duration
	Logf            LogFunc
}
