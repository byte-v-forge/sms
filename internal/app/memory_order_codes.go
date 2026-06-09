package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *MemoryOrderStore) RecordCode(ctx context.Context, order core.Order, code core.SMSCode, _ ...eventoutbox.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = cloneOrder(order)
	s.codes[order.ID] = append(s.codes[order.ID], core.OrderCode{OrderID: order.ID, Code: cloneSMSCode(code)})
	s.pruneLocked(s.clock.Now())
	return nil
}
