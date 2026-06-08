package app

import (
	"context"
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
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

func (s *MemoryOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, code := range s.codes[orderID] {
		if code.Code.SecretRef.GetSecretId() == secretID {
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryOrderStore) ListCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limitPerOrder = pagex.NormalizeLimit(limitPerOrder, defaultOrderCodeLimit, maxOrderCodeLimit)
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []core.OrderCode{}
	for _, orderID := range orderIDs {
		codes := append([]core.OrderCode{}, s.codes[orderID]...)
		sort.Slice(codes, func(i, j int) bool {
			return codes[i].Code.ReceivedAt.After(codes[j].Code.ReceivedAt)
		})
		if len(codes) > limitPerOrder {
			codes = codes[:limitPerOrder]
		}
		for _, code := range codes {
			out = append(out, core.OrderCode{OrderID: code.OrderID, Code: cloneSMSCode(code.Code)})
		}
	}
	return out, nil
}
