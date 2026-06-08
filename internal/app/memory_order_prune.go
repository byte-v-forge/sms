package app

import "time"

func (s *MemoryOrderStore) pruneLocked(now time.Time) {
	if len(s.orders) == 0 {
		return
	}
	s.pruneFinalOrdersLocked(now)
	if len(s.orders) <= memoryOrderMaxEntries {
		return
	}
	s.pruneOverflowOrdersLocked()
}
