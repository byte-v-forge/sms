package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func PublicOrder(order core.Order) *smsv1.SmsOrder {
	return &smsv1.SmsOrder{
		OrderId:   order.ID,
		RequestId: order.RequestID,
		Target: &smsv1.SmsTarget{
			ApplicationKey:     order.Target.ApplicationKey,
			CountryIso2:        order.Target.CountryISO2,
			CountryCallingCode: order.Target.CountryCallingCode,
		},
		PhoneNumber: &smsv1.PhoneNumber{
			E164Number:         order.PhoneNumber.E164,
			NationalNumber:     order.PhoneNumber.NationalNumber,
			CountryIso2:        order.PhoneNumber.CountryISO2,
			CountryCallingCode: order.PhoneNumber.CountryCallingCode,
		},
		Status:                   PublicOrderStatus(order.Status),
		Price:                    PublicMoney(order.Price),
		AcquiredAt:               PublicTime(order.AcquiredAt),
		ExpiresAt:                PublicTime(order.ExpiresAt),
		UpdatedAt:                PublicTime(order.UpdatedAt),
		LastError:                PublicError(order.LastError),
		CanRequestAdditionalCode: order.CanRequestAdditionalCode,
		CancelAllowedAt:          PublicTime(order.CancelAllowedAt),
	}
}
