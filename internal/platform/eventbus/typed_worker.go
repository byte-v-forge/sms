package eventbus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

type MessageAction string

const (
	MessageActionAck  MessageAction = "ack"
	MessageActionNak  MessageAction = "nak"
	MessageActionTerm MessageAction = "term"
)

type HandlerResult struct {
	Action MessageAction
	Delay  time.Duration
	Label  string
}

func AckResult(label string) HandlerResult {
	return HandlerResult{Action: MessageActionAck, Label: label}
}

func NakResult(delay time.Duration, label string) HandlerResult {
	return HandlerResult{Action: MessageActionNak, Delay: delay, Label: label}
}

func TermResult(label string) HandlerResult {
	return HandlerResult{Action: MessageActionTerm, Label: label}
}

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

func RunTypedConsumerWorker[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T]) error {
	if cfg.Consumer == nil || cfg.NewMessage == nil || cfg.Handler == nil {
		return nil
	}
	cfg = normalizeTypedConsumerWorkerConfig(cfg)
	if err := cfg.Expected.ValidateEvent(cfg.NewMessage()); err != nil {
		return fmt.Errorf("configure %s expected event: %w", cfg.Name, err)
	}
	return RunConsumerWorker(ctx, ConsumerWorkerConfig{
		Name:            cfg.Name,
		Consumer:        cfg.Consumer,
		Batch:           cfg.Batch,
		FetchErrorDelay: cfg.FetchErrorDelay,
		Logf:            cfg.Logf,
		Handler: func(ctx context.Context, received ReceivedMessage) {
			handleTypedMessage(ctx, cfg, received)
		},
	})
}

func handleTypedMessage[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T], received ReceivedMessage) {
	if err := cfg.Expected.ValidateReceived(received); err != nil {
		cfg.Logf("validate %s envelope failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return
	}
	message := cfg.NewMessage()
	if err := UnmarshalPayload(received, message); err != nil {
		cfg.Logf("decode %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
		TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
		return
	}
	if cfg.Validate != nil {
		if err := cfg.Validate(message); err != nil {
			cfg.Logf("validate %s failed event_id=%s: %v", cfg.Name, EventID(received), err)
			TermMessage(ctx, received, cfg.MalformedLabel, cfg.Logf)
			return
		}
	}
	applyHandlerResult(ctx, received, cfg.Handler(ctx, message, received), cfg.Logf)
}

func applyHandlerResult(ctx context.Context, message ReceivedMessage, result HandlerResult, logf LogFunc) {
	label := strings.TrimSpace(result.Label)
	switch result.Action {
	case MessageActionNak:
		if label == "" {
			label = "nak event"
		}
		NakMessageDelay(ctx, message, result.Delay, label, logf)
	case MessageActionTerm:
		if label == "" {
			label = "terminate event"
		}
		TermMessage(ctx, message, label, logf)
	default:
		if label == "" {
			label = "ack event"
		}
		AckMessage(ctx, message, label, logf)
	}
}

func normalizeTypedConsumerWorkerConfig[T proto.Message](cfg TypedConsumerWorkerConfig[T]) TypedConsumerWorkerConfig[T] {
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = "typed event consumer"
	}
	cfg.MalformedLabel = strings.TrimSpace(cfg.MalformedLabel)
	if cfg.MalformedLabel == "" {
		cfg.MalformedLabel = fmt.Sprintf("terminate malformed %s", cfg.Name)
	}
	cfg.Logf = logger(cfg.Logf)
	return cfg
}
