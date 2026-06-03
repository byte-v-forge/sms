package core

import (
	"time"

	commonv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/common/v1"
)

type SMSCode struct {
	Value       string
	MessageText string
	ReceivedAt  time.Time
	SecretRef   *commonv1.SecretRef
}

type OrderCode struct {
	OrderID string
	Code    SMSCode
}

type Order struct {
	ID                       string
	RequestID                string
	ProviderKey              string
	UpstreamOrderID          string
	Target                   Target
	PhoneNumber              PhoneNumber
	Status                   OrderStatus
	Price                    Money
	AcquiredAt               time.Time
	ExpiresAt                time.Time
	UpdatedAt                time.Time
	CancelAllowedAt          time.Time
	CanRequestAdditionalCode bool
	LastError                *Error
}

func (a Order) IsExpired(now time.Time) bool {
	return !a.ExpiresAt.IsZero() && !a.Status.IsFinal() && !now.Before(a.ExpiresAt)
}
