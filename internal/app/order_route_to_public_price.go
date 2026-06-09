package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func applyPublicAcquirePriceParams(params *smsv1.SmsNumberAcquireParams, route core.Route) {
	if moneyIsSet(route.MinPrice) {
		params.MinPrice = PublicMoney(route.MinPrice)
	}
	if moneyIsSet(route.MaxPrice) {
		params.MaxPrice = PublicMoney(route.MaxPrice)
	}
}
