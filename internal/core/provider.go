package core

import "time"

type ProviderPolicy struct {
	OrderTTL              time.Duration
	PollInterval          time.Duration
	CancelAllowedAfter    time.Duration
	EarlyCancelRetryAfter time.Duration
	CancelAllowedUntil    time.Duration
}

func (p ProviderPolicy) WithDefaults() ProviderPolicy {
	if p.OrderTTL <= 0 {
		p.OrderTTL = 20 * time.Minute
	}
	if p.PollInterval <= 0 {
		p.PollInterval = 5 * time.Second
	}
	return p
}
