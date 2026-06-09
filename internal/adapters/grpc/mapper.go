package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

func toProtoOrder(order core.Order) *smsv1.SmsOrder { return app.PublicOrder(order) }
func toProtoCode(code *core.SMSCode) *smsv1.SmsCode { return app.PublicCode(code) }
func toProtoMoney(money core.Money) *smsv1.DecimalMoney {
	return app.PublicMoney(money)
}
func toProtoError(err error) *smsv1.SmsError { return app.PublicError(err) }
func toProtoStatus(status core.OrderStatus) smsv1.SmsOrderStatus {
	return app.PublicOrderStatus(status)
}
