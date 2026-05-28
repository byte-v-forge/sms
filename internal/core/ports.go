package core

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) string
}

type Provider interface {
	Key() string
	Policy() ProviderPolicy
	AcquireNumber(ctx context.Context, request ProviderAcquireRequest) (ProviderOrder, error)
	GetStatus(ctx context.Context, upstreamOrderID string) (ProviderCodeResult, error)
	SetStatus(ctx context.Context, upstreamOrderID string, action ProviderAction) error
	GetBalance(ctx context.Context) (Money, error)
}
