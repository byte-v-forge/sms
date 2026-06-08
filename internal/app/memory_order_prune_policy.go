package app

import "time"

const (
	memoryOrderMaxEntries     = 1000
	memoryOrderFinalRetention = 2 * time.Hour
)

type memoryOrderEntry struct {
	id        string
	final     bool
	updatedAt time.Time
}
