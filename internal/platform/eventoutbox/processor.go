package eventoutbox

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/jackc/pgx/v5"
)

type PgxBeginner interface {
	Begin(context.Context) (pgx.Tx, error)
}

type PgxProcessor struct {
	Beginner       PgxBeginner
	Table          string
	Publisher      eventbus.Publisher
	PublishOptions PublishOptions
}

func NewPgxProcessor(beginner PgxBeginner, table string, publisher eventbus.Publisher) *PgxProcessor {
	return &PgxProcessor{Beginner: beginner, Table: table, Publisher: publisher}
}

func (p *PgxProcessor) PublishPending(ctx context.Context, batch int) (int, error) {
	if p == nil || p.Beginner == nil || p.Publisher == nil {
		return 0, nil
	}
	if batch <= 0 {
		batch = DefaultBatch
	}
	tx, err := p.Beginner.Begin(ctx)
	if err != nil {
		return 0, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	rows, err := ClaimPendingPgx(ctx, tx, p.Table, batch, optionUnix(p.PublishOptions))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	updates, err := NewPgxUpdates(tx, p.Table)
	if err != nil {
		return 0, err
	}
	published, err := PublishRows(ctx, p.Publisher, rows, updates, p.PublishOptions)
	if err != nil {
		return published, err
	}
	if err := tx.Commit(ctx); err != nil {
		return published, err
	}
	committed = true
	return published, nil
}

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

func optionUnix(options PublishOptions) int64 {
	if options.Now == nil {
		return 0
	}
	return options.Now().Unix()
}
