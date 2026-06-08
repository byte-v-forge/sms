package app

import (
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PublicCode(code *core.SMSCode) *smsv1.SmsCode {
	if code == nil {
		return nil
	}
	return &smsv1.SmsCode{
		SecretRef:  secretref.Clone(code.SecretRef, "sms", "sms_otp"),
		ReceivedAt: PublicTime(code.ReceivedAt),
	}
}

func PublicMoney(money core.Money) *smsv1.DecimalMoney {
	if money.CurrencyCode == "" && money.AmountDecimal == "" {
		return nil
	}
	return &smsv1.DecimalMoney{CurrencyCode: money.CurrencyCode, AmountDecimal: money.AmountDecimal}
}

func PublicTime(value time.Time) *timestamppb.Timestamp {
	if value.IsZero() {
		return nil
	}
	return timestamppb.New(value)
}
