package app

import "context"

func (s *MemoryTTLStringStore) Load(ctx context.Context, key string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	normalized := s.key(key)
	if normalized == "" {
		return "", false, nil
	}
	s.mu.RLock()
	item, ok := s.values[normalized]
	s.mu.RUnlock()
	if !ok {
		return "", false, nil
	}
	if !item.expiresAt.IsZero() && !s.clock.Now().Before(item.expiresAt) {
		s.mu.Lock()
		delete(s.values, normalized)
		s.mu.Unlock()
		return "", false, nil
	}
	return item.value, true, nil
}
