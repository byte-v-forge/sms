package eventoutbox

import (
	"context"

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
