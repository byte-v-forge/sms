package core

type ProviderAction string

const (
	ActionMarkMessageSent   ProviderAction = "mark_message_sent"
	ActionRequestAdditional ProviderAction = "request_additional_code"
	ActionCompleteOrder     ProviderAction = "complete_order"
	ActionCancelOrder       ProviderAction = "cancel_order"
)
