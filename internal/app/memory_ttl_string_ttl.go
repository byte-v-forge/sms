package app

import "time"

func (s *MemoryTTLStringStore) DefaultTTL() time.Duration {
	if s == nil {
		return 0
	}
	return s.ttl
}
