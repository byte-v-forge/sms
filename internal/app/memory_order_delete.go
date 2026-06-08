package app

func (s *MemoryOrderStore) deleteOrderLocked(orderID string) {
	delete(s.orders, orderID)
	delete(s.codes, orderID)
}
