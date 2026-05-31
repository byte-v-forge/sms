package core

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
