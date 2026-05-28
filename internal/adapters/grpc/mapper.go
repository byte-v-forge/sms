package grpcadapter

import (
	"time"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoOrder(order core.Order) *smsv1.SmsOrder {
	return app.PublicOrder(order)
}
func toProtoCode(code *core.SMSCode) *smsv1.SmsCode     { return app.PublicCode(code) }
func toProtoMoney(money core.Money) *smsv1.DecimalMoney { return app.PublicMoney(money) }
func toProtoError(err error) *smsv1.SmsError            { return app.PublicError(err) }
func toProtoStatus(status core.OrderStatus) smsv1.SmsOrderStatus {
	return app.PublicOrderStatus(status)
}
func toProtoTime(value time.Time) *timestamppb.Timestamp { return app.PublicTime(value) }

func fromProtoMoney(value *smsv1.DecimalMoney) core.Money {
	if value == nil {
		return core.Money{}
	}
	return core.Money{CurrencyCode: value.GetCurrencyCode(), AmountDecimal: value.GetAmountDecimal()}
}

func fromProtoTarget(target *smsv1.SmsTarget) core.Target {
	if target == nil {
		return core.Target{}
	}
	return core.Target{
		ApplicationKey:     target.GetApplicationKey(),
		CountryISO2:        target.GetCountryIso2(),
		CountryCallingCode: target.GetCountryCallingCode(),
	}
}

func protoDuration(value *durationpb.Duration) time.Duration {
	if value == nil {
		return 0
	}
	return value.AsDuration()
}
