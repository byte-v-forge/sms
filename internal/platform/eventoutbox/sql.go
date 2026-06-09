package eventoutbox

import "errors"

var (
	ErrInvalidTableName = errors.New("event outbox table name is invalid")
	ErrNilDB            = errors.New("event outbox database handle is nil")
)
