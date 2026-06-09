package app

import "github.com/byte-v-forge/sms/internal/core"

func (s *MemoryOrderStore) filteredOrdersLocked(includeFinal bool) []core.Order {
	orders := make([]core.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if !includeFinal && order.Status.IsFinal() {
			continue
		}
		orders = append(orders, cloneOrder(order))
	}
	return orders
}
