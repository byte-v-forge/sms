package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
)

func (s *ProviderAdminService) ListOrderCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if s.orderDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms code history is not available", false)
	}
	ids := uniqueNonEmpty(orderIDs)
	if len(ids) == 0 {
		return []core.OrderCode{}, nil
	}
	if len(ids) > 200 {
		return nil, core.NewError(core.CodeValidationFailed, "too many order_ids", false)
	}
	return s.orderDB.ListCodes(ctx, ids, pagex.NormalizeLimit(limitPerOrder, defaultOrderCodeLimit, maxOrderCodeLimit))
}
