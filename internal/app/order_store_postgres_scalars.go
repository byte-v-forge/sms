package app

import (
	"database/sql"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func timeOrNil(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func nullableTime(value sql.NullTime) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func errorCode(err *core.Error) string {
	if err == nil {
		return ""
	}
	return string(err.Code)
}

func errorFromCode(code string) *core.Error {
	if code == "" {
		return nil
	}
	return &core.Error{Code: core.ErrorCode(code)}
}
