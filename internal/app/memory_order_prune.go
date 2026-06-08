package app

import (
	"sort"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

const (
	memoryOrderMaxEntries     = 1000
	memoryOrderFinalRetention = 2 * time.Hour
)

type memoryOrderEntry struct {
	id        string
	final     bool
	updatedAt time.Time
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
