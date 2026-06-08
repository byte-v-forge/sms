package grpcadapter

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func toOrderView(order core.Order) *smsinternalv1.SmsOrderView {
	if order.ID == "" {
		return nil
	}
	return &smsinternalv1.SmsOrderView{
		Order:       toProtoOrder(order),
		ProviderKey: order.ProviderKey,
	}
}

func toOrderCodeView(code core.OrderCode) *smsinternalv1.SmsOrderCodeView {
	return &smsinternalv1.SmsOrderCodeView{
		OrderId: code.OrderID,
		Code:    toProtoCode(&code.Code),
	}
}

func toProviderError(err error) *smsinternalv1.ProviderError {
	if err == nil {
		return nil
	}
	return &smsinternalv1.ProviderError{PublicError: toProtoError(err)}
}
