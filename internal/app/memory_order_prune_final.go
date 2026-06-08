package app

import "time"

func (s *MemoryOrderStore) pruneFinalOrdersLocked(now time.Time) {
	for orderID, order := range s.orders {
		if order.Status.IsFinal() && memoryOrderAge(order, now) >= memoryOrderFinalRetention {
			s.deleteOrderLocked(orderID)
		}
	}
}
