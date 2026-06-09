package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
)

func (s *MemoryOrderStore) ListCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limitPerOrder = pagex.NormalizeLimit(limitPerOrder, defaultOrderCodeLimit, maxOrderCodeLimit)
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []core.OrderCode{}
	for _, orderID := range orderIDs {
		out = append(out, limitedMemoryOrderCodes(s.codes[orderID], limitPerOrder)...)
	}
	return out, nil
}
