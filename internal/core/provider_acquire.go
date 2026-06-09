package core

import "time"

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
