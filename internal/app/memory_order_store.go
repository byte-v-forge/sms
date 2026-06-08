package app

import (
	"context"
	"sort"
	"sync"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
	"google.golang.org/protobuf/proto"
)

type MemoryOrderStore struct {
	mu     sync.RWMutex
	orders map[string]core.Order
	codes  map[string][]core.OrderCode
	clock  core.Clock
}

const (
	memoryOrderMaxEntries     = 1000
	memoryOrderFinalRetention = 2 * time.Hour
)

func NewMemoryOrderStore() *MemoryOrderStore {
	return &MemoryOrderStore{orders: map[string]core.Order{}, codes: map[string][]core.OrderCode{}, clock: SystemClock{}}
}

func (s *MemoryOrderStore) Save(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *MemoryOrderStore) Update(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

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

func (s *MemoryOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	if err := ctx.Err(); err != nil {
		return core.Order{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderID]
	if !ok {
		return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false)
	}
	return cloneOrder(order), nil
}

func (s *MemoryOrderStore) List(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limit = pagex.NormalizeLimit(limit, pagex.DefaultLimit, pagex.MaxLimit)
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]core.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if !includeFinal && order.Status.IsFinal() {
			continue
		}
		orders = append(orders, cloneOrder(order))
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].UpdatedAt.Equal(orders[j].UpdatedAt) {
			return orders[i].ID < orders[j].ID
		}
		return orders[i].UpdatedAt.After(orders[j].UpdatedAt)
	})
	if len(orders) > limit {
		return orders[:limit], nil
	}
	return orders, nil
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

func (s *MemoryOrderStore) save(ctx context.Context, order core.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = cloneOrder(order)
	s.pruneLocked(s.clock.Now())
	return nil
}

func (s *MemoryOrderStore) pruneLocked(now time.Time) {
	if len(s.orders) == 0 {
		return
	}
	for orderID, order := range s.orders {
		if order.Status.IsFinal() && memoryOrderAge(order, now) >= memoryOrderFinalRetention {
			s.deleteOrderLocked(orderID)
		}
	}
	if len(s.orders) <= memoryOrderMaxEntries {
		return
	}
	entries := make([]memoryOrderEntry, 0, len(s.orders))
	for orderID, order := range s.orders {
		entries = append(entries, memoryOrderEntry{id: orderID, final: order.Status.IsFinal(), updatedAt: memoryOrderUpdatedAt(order)})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].final != entries[j].final {
			return entries[i].final
		}
		if entries[i].updatedAt.Equal(entries[j].updatedAt) {
			return entries[i].id < entries[j].id
		}
		return entries[i].updatedAt.Before(entries[j].updatedAt)
	})
	for _, entry := range entries {
		if len(s.orders) <= memoryOrderMaxEntries {
			return
		}
		s.deleteOrderLocked(entry.id)
	}
}

type memoryOrderEntry struct {
	id        string
	final     bool
	updatedAt time.Time
}

func (s *MemoryOrderStore) deleteOrderLocked(orderID string) {
	delete(s.orders, orderID)
	delete(s.codes, orderID)
}

func memoryOrderAge(order core.Order, now time.Time) time.Duration {
	updatedAt := memoryOrderUpdatedAt(order)
	if updatedAt.IsZero() || now.Before(updatedAt) {
		return 0
	}
	return now.Sub(updatedAt)
}

func memoryOrderUpdatedAt(order core.Order) time.Time {
	for _, value := range []time.Time{order.UpdatedAt, order.ExpiresAt, order.AcquiredAt} {
		if !value.IsZero() {
			return value
		}
	}
	return time.Time{}
}

func cloneOrder(order core.Order) core.Order {
	order.LastError = cloneCoreError(order.LastError)
	return order
}

func cloneCoreError(err *core.Error) *core.Error {
	if err == nil {
		return nil
	}
	out := *err
	return &out
}

func cloneSMSCode(code core.SMSCode) core.SMSCode {
	code.SecretRef = cloneSecretRef(code.SecretRef)
	return code
}

func cloneSecretRef(ref *commonv1.SecretRef) *commonv1.SecretRef {
	if ref == nil {
		return nil
	}
	return proto.Clone(ref).(*commonv1.SecretRef)
}
