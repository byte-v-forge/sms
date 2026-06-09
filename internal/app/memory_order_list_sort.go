package app

import (
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
)

func sortOrdersByUpdatedAt(orders []core.Order) {
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].UpdatedAt.Equal(orders[j].UpdatedAt) {
			return orders[i].ID < orders[j].ID
		}
		return orders[i].UpdatedAt.After(orders[j].UpdatedAt)
	})
}
