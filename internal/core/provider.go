package core

import "time"

type ProviderAction string

const (
	ActionMarkMessageSent   ProviderAction = "mark_message_sent"
	ActionRequestAdditional ProviderAction = "request_additional_code"
	ActionCompleteOrder     ProviderAction = "complete_order"
	ActionCancelOrder       ProviderAction = "cancel_order"
)

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

type AcquireNumberCommand struct {
	RequestID     string
	AcquireParams Route
	LeaseDuration time.Duration
}

type ProviderAcquireRequest struct {
	RequestID     string
	Route         Route
	Target        Target
	LeaseDuration time.Duration
}

type ProviderOrder struct {
	UpstreamOrderID          string
	PhoneNumber              PhoneNumber
	Price                    Money
	AcquiredAt               time.Time
	ExpiresAt                time.Time
	CanRequestAdditionalCode bool
}

type ProviderCodeResult struct {
	Status      OrderStatus
	Code        string
	MessageText string
	ReceivedAt  time.Time
}
