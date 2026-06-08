package app

import (
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func recommendationPriceRange(policy *smsv1.SmsRoutePolicy) (core.Money, core.Money, error) {
	minPrice, err := normalizeRecommendationPrice(moneyFromProto(policy.GetMinPrice()), "min_price")
	if err != nil {
		return core.Money{}, core.Money{}, err
	}
	maxPrice, err := normalizeRecommendationPrice(moneyFromProto(policy.GetMaxPrice()), "max_price")
	if err != nil {
		return core.Money{}, core.Money{}, err
	}
	if moneyIsSet(minPrice) && moneyIsSet(maxPrice) {
		if minPrice.CurrencyCode != "" && maxPrice.CurrencyCode != "" && minPrice.CurrencyCode != maxPrice.CurrencyCode {
			return core.Money{}, core.Money{}, core.NewError(core.CodeValidationFailed, "sms route recommendation price currency mismatch", false)
		}
		minAmount, _ := parseDecimalAmount(minPrice.AmountDecimal)
		maxAmount, _ := parseDecimalAmount(maxPrice.AmountDecimal)
		if minAmount.Cmp(maxAmount) > 0 {
			return core.Money{}, core.Money{}, core.NewError(core.CodeValidationFailed, "sms route recommendation min_price exceeds max_price", false)
		}
	}
	return minPrice, maxPrice, nil
}

func normalizeRecommendationPrice(price core.Money, field string) (core.Money, error) {
	if strings.TrimSpace(price.AmountDecimal) == "" {
		return core.Money{}, nil
	}
	amount, ok := parseDecimalAmount(price.AmountDecimal)
	if !ok || amount.Sign() < 0 {
		return core.Money{}, core.NewError(core.CodeValidationFailed, "sms route recommendation "+field+" is invalid", false)
	}
	price.CurrencyCode = strings.ToUpper(strings.TrimSpace(price.CurrencyCode))
	return price, nil
}
