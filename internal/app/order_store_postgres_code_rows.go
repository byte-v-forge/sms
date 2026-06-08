package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
	"github.com/jackc/pgx/v5"
)

func scanPostgresOrderCodes(rows pgx.Rows) ([]core.OrderCode, error) {
	out := []core.OrderCode{}
	for rows.Next() {
		item, err := scanPostgresOrderCode(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func scanPostgresOrderCode(rows pgx.Rows) (core.OrderCode, error) {
	var item core.OrderCode
	var secretID string
	var expiresAt time.Time
	if err := rows.Scan(&item.OrderID, &secretID, &item.Code.MessageText, &item.Code.ReceivedAt, &expiresAt); err != nil {
		return core.OrderCode{}, err
	}
	item.Code.SecretRef = secretref.New("sms", "sms_otp", secretID, expiresAt)
	return item, nil
}
