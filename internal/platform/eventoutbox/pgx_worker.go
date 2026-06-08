package eventoutbox

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

type PgxWorkerConfig struct {
	Name           string
	Beginner       PgxBeginner
	Table          string
	Publisher      eventbus.Publisher
	Batch          int
	Interval       time.Duration
	ActiveInterval time.Duration
	Logf           func(string, ...any)
	PublishOptions PublishOptions
}

func RunPgxWorker(ctx context.Context, cfg PgxWorkerConfig) error {
	if cfg.Beginner == nil || cfg.Publisher == nil {
		return nil
	}
	return RunWorker(ctx, WorkerConfig{
		Name:           cfg.Name,
		Processor:      &PgxProcessor{Beginner: cfg.Beginner, Table: cfg.Table, Publisher: cfg.Publisher, PublishOptions: cfg.PublishOptions},
		Batch:          cfg.Batch,
		Interval:       cfg.Interval,
		ActiveInterval: cfg.ActiveInterval,
		Logf:           cfg.Logf,
	})
}
