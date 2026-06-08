package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
)

func (s *ProviderAdminService) ListOrders(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if s.orderDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms order list is not available", false)
	}
	return s.orderDB.List(ctx, includeFinal, pagex.NormalizeLimit(limit, pagex.DefaultLimit, pagex.MaxLimit))
}
