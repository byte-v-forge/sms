package app

import "strings"

func (s *MemoryTTLStringStore) key(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if s.keyspace == "" {
		return value
	}
	return s.keyspace + ":" + value
}
