package core

import "time"

type ProviderCodeResult struct {
	Status      OrderStatus
	Code        string
	MessageText string
	ReceivedAt  time.Time
}
