package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func moneyFromProto(value *smsv1.DecimalMoney) core.Money {
	if value == nil {
		return core.Money{}
	}
	return core.Money{CurrencyCode: value.GetCurrencyCode(), AmountDecimal: value.GetAmountDecimal()}
}
