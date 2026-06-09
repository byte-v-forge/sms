package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

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
