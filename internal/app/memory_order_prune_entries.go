package app

import (
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
)

func sortedMemoryOrderEntries(orders map[string]core.Order) []memoryOrderEntry {
	entries := make([]memoryOrderEntry, 0, len(orders))
	for orderID, order := range orders {
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
	return entries
}
