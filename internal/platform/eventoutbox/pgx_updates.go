package eventoutbox

func NewPgxUpdates(tx PgxTx, table string) (Updates, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	return pgxUpdates{tx: tx, table: tableName}, nil
}

type pgxUpdates struct {
	tx    PgxTx
	table string
}
