package eventbus

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

const DefaultFetchErrorDelay = time.Second

type LogFunc func(string, ...any)

type MessageHandler func(context.Context, ReceivedMessage)

type ConsumerWorkerConfig struct {
	Name            string
	Consumer        Consumer
	Handler         MessageHandler
	Batch           int
	FetchErrorDelay time.Duration
	Logf            LogFunc
}

func RunConsumerWorker(ctx context.Context, cfg ConsumerWorkerConfig) error {
	if cfg.Consumer == nil || cfg.Handler == nil {
		return nil
	}
	cfg = normalizeConsumerWorkerConfig(cfg)
	for ctx.Err() == nil {
		messages, err := cfg.Consumer.Fetch(ctx, cfg.Batch)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			cfg.Logf("fetch %s failed: %v", cfg.Name, err)
			if err := timex.Sleep(ctx, cfg.FetchErrorDelay); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				return err
			}
			continue
		}
		for _, message := range messages {
			cfg.Handler(ctx, message)
		}
	}
	return nil
}

func Ack(ctx context.Context, action func(context.Context) error, label string, logf LogFunc) {
	if action == nil {
		return
	}
	if err := action(ctx); err != nil && ctx.Err() == nil {
		logger(logf)("%s failed: %v", label, err)
	}
}

func AckMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	Ack(ctx, message.Ack, label, logf)
}

func NakMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	Ack(ctx, message.Nak, label, logf)
}

func NakMessageDelay(ctx context.Context, message ReceivedMessage, delay time.Duration, label string, logf LogFunc) {
	if delay > 0 && message.NakDelay != nil {
		Ack(ctx, func(nakCtx context.Context) error { return message.NakDelay(nakCtx, delay) }, label, logf)
		return
	}
	NakMessage(ctx, message, label, logf)
}

func TermMessage(ctx context.Context, message ReceivedMessage, label string, logf LogFunc) {
	if message.DeadLetter != nil {
		Ack(ctx, func(deadLetterCtx context.Context) error {
			return message.DeadLetter(deadLetterCtx, label)
		}, "publish dead letter for "+label, logf)
	}
	Ack(ctx, message.Term, label, logf)
}

func EventID(message ReceivedMessage) string {
	if message.Envelope == nil || message.Envelope.GetMetadata() == nil {
		return ""
	}
	return message.Envelope.GetMetadata().GetId()
}

func normalizeConsumerWorkerConfig(cfg ConsumerWorkerConfig) ConsumerWorkerConfig {
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = "event consumer"
	}
	if cfg.FetchErrorDelay <= 0 {
		cfg.FetchErrorDelay = DefaultFetchErrorDelay
	}
	cfg.Logf = logger(cfg.Logf)
	return cfg
}

func logger(logf LogFunc) LogFunc {
	if logf != nil {
		return logf
	}
	return log.Printf
}
