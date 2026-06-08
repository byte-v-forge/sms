package app

func (s *MemoryOrderStore) pruneOverflowOrdersLocked() {
	entries := sortedMemoryOrderEntries(s.orders)
	for _, entry := range entries {
		if len(s.orders) <= memoryOrderMaxEntries {
			return
		}
		s.deleteOrderLocked(entry.id)
	}
}
