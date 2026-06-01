package app

import (
	"sort"
	"strings"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/common-lib/pagex"
	"github.com/byte-v-forge/sms/internal/core"
)

func sortedProviderFilterKeys(providerFilter map[string]struct{}) []string {
	keys := make([]string, 0, len(providerFilter))
	for providerKey := range providerFilter {
		keys = append(keys, providerKey)
	}
	sort.Strings(keys)
	return keys
}

func normalizeRecommendationTarget(target core.Target) core.Target {
	target.ApplicationKey = strings.TrimSpace(target.ApplicationKey)
	target.CountryISO2 = strings.ToUpper(strings.TrimSpace(target.CountryISO2))
	target.CountryCallingCode = strings.TrimPrefix(strings.TrimSpace(target.CountryCallingCode), "+")
	return target
}

func validateRecommendationTarget(target core.Target) error {
	if target.ApplicationKey == "" {
		return core.NewError(core.CodeValidationFailed, "sms route recommendation target application_key is required", false)
	}
	if target.CountryISO2 == "" && target.CountryCallingCode == "" {
		return core.NewError(core.CodeValidationFailed, "sms route recommendation target country is required", false)
	}
	return nil
}

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

func recommendationStrategy(policy *smsv1.SmsRoutePolicy) smsv1.SmsRouteStrategy {
	switch policy.GetStrategy() {
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE
	default:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE
	}
}

func recommendationLimit(policy *smsv1.SmsRoutePolicy) int {
	return pagex.NormalizeLimit(int(policy.GetLimit()), defaultRouteRecommendationLimit, maxRouteRecommendationLimit)
}

func normalizedProviderFilter(providerKeys []string) map[string]struct{} {
	if len(providerKeys) == 0 {
		return nil
	}
	filter := make(map[string]struct{}, len(providerKeys))
	for _, providerKey := range providerKeys {
		key := normalizeProviderKey(providerKey)
		if key == "" {
			continue
		}
		filter[key] = struct{}{}
	}
	if len(filter) == 0 {
		return nil
	}
	return filter
}
