package app

import "context"

func (s *PostgresOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
SELECT EXISTS (
  SELECT 1
  FROM sms_order_codes
  WHERE order_id = $1 AND code_secret_id = $2
)
`, orderID, secretID).Scan(&exists)
	return exists, err
}
