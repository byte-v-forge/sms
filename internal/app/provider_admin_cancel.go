package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *ProviderAdminService) CancelOrder(ctx context.Context, orderID string, requestID string) (core.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return core.Order{}, core.NewError(core.CodeValidationFailed, "order_id is required", false)
	}
	if strings.TrimSpace(requestID) == "" {
		requestID = RandomIDGenerator{}.NewID("req_")
	}
	return s.orders.CancelOrder(ctx, orderID, requestID)
}
