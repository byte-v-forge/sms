package eventoutbox

import "context"

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
