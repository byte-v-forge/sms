package core

import (
	"fmt"
	"time"
)

type Money struct {
	CurrencyCode  string
	AmountDecimal string
}

type PhoneNumber struct {
	E164               string
	NationalNumber     string
	CountryISO2        string
	CountryCallingCode string
}

type Target struct {
	ApplicationKey     string
	CountryISO2        string
	CountryCallingCode string
}

type OrderStatus string

const (
	StatusAcquireRequested        OrderStatus = "acquire_requested"
	StatusPendingCode             OrderStatus = "pending_code"
	StatusMessageSent             OrderStatus = "message_sent"
	StatusCodeReceived            OrderStatus = "code_received"
	StatusAdditionalCodeRequested OrderStatus = "additional_code_requested"
	StatusCompleted               OrderStatus = "completed"
	StatusCanceled                OrderStatus = "canceled"
	StatusExpired                 OrderStatus = "expired"
	StatusFailed                  OrderStatus = "failed"
)

func (s OrderStatus) IsFinal() bool {
	switch s {
	case StatusCompleted, StatusCanceled, StatusExpired, StatusFailed:
		return true
	default:
		return false
	}
}

func (s OrderStatus) HasProviderLease() bool {
	switch s {
	case StatusPendingCode, StatusMessageSent, StatusCodeReceived, StatusAdditionalCodeRequested, StatusCompleted, StatusCanceled, StatusExpired, StatusFailed:
		return true
	default:
		return false
	}
}

type ErrorCode string

const (
	CodeValidationFailed      ErrorCode = "validation_failed"
	CodeRouteNotFound         ErrorCode = "route_not_found"
	CodeOrderNotFound         ErrorCode = "order_not_found"
	CodeOrderAlreadyFinalized ErrorCode = "order_already_finalized"
	CodeNoNumberAvailable     ErrorCode = "no_number_available"
	CodeRateLimited           ErrorCode = "rate_limited"
	CodeSupplyUnavailable     ErrorCode = "supply_unavailable"
	CodeUpstreamRejected      ErrorCode = "upstream_rejected"
	CodeTimeout               ErrorCode = "timeout"
	CodeUnsupportedOperation  ErrorCode = "unsupported_operation"
	CodeOrderExpired          ErrorCode = "order_expired"
	CodeCancelNotAllowed      ErrorCode = "cancel_not_allowed"
	CodeInsufficientBalance   ErrorCode = "insufficient_balance"
	CodeInternal              ErrorCode = "internal"
)

type Error struct {
	Code      ErrorCode
	Message   string
	Retryable bool
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, message string, retryable bool) *Error {
	return &Error{Code: code, Message: message, Retryable: retryable}
}

type SMSCode struct {
	Value       string
	MessageText string
	ReceivedAt  time.Time
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

type Route struct {
	ProviderKey        string
	ApplicationKey     string
	UpstreamServiceKey string
	CountryISO2        string
	CountryCallingCode string
	ProviderCountryID  string
	UpstreamProviderID string
}

type RouteOfferQuery struct {
	ApplicationKey     string
	CountryISO2        string
	CountryCallingCode string
	ProviderKey        string
}

type RouteOffer struct {
	ProviderKey             string
	ProviderDisplayName     string
	UpstreamProviderID      string
	UpstreamProviderName    string
	ApplicationKey          string
	ApplicationName         string
	CountryISO2             string
	CountryName             string
	CountryCallingCode      string
	Price                   Money
	AvailableCount          int
	SupportsCancel          bool
	SupportsAdditionalCode  bool
	RequiresMarkMessageSent bool
	ObservedAt              time.Time
	Route                   Route
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
