package eventbusadapter

import (
	"errors"
	"strings"
)

func validateOrderID(orderID string) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.New("order_id is required")
	}
	return nil
}
