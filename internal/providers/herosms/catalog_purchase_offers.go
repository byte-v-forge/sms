package herosms

import (
	"sort"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func purchaseInfoPriceOffers(serviceKey, countryID string, info purchaseInfo) []PriceOffer {
	offers := make([]PriceOffer, 0)
	for _, operator := range info.Operators {
		offers = append(offers, purchaseOperatorPriceOffers(serviceKey, countryID, info, operator)...)
	}
	sortPurchaseInfoPriceOffers(offers)
	return uniqueHeroSMSOffers(offers)
}

func purchaseOperatorPriceOffers(serviceKey, countryID string, info purchaseInfo, operator purchaseOperator) []PriceOffer {
	if len(operator.FreePriceOffers) > 0 {
		return purchaseOperatorFreePriceOffers(serviceKey, countryID, operator)
	}
	return purchaseOperatorDefaultPriceOffer(serviceKey, countryID, info, operator)
}

func purchaseOperatorFreePriceOffers(serviceKey, countryID string, operator purchaseOperator) []PriceOffer {
	offers := make([]PriceOffer, 0, len(operator.FreePriceOffers))
	for price, count := range operator.FreePriceOffers {
		if count <= 0 || strings.TrimSpace(price) == "" {
			continue
		}
		offers = append(offers, priceOffer(serviceKey, countryID, operator, price, count))
	}
	return offers
}

func purchaseOperatorDefaultPriceOffer(serviceKey, countryID string, info purchaseInfo, operator purchaseOperator) []PriceOffer {
	price := jsonx.FirstScalar(info.UserPrice)
	if price == "" || operator.ActivationsCount <= 0 {
		return nil
	}
	return []PriceOffer{priceOffer(serviceKey, countryID, operator, price, operator.ActivationsCount)}
}

func priceOffer(serviceKey, countryID string, operator purchaseOperator, price string, count int) PriceOffer {
	operatorKey := heroSMSOperator(operator.Name)
	return PriceOffer{
		CountryID:          strings.TrimSpace(countryID),
		UpstreamServiceKey: normalizeHeroSMSServiceKey(serviceKey),
		Operator:           operatorKey,
		OperatorName:       strings.TrimSpace(operator.LocalName),
		Price:              core.Money{CurrencyCode: "USD", AmountDecimal: strings.TrimSpace(price)},
		AvailableCount:     count,
	}
}

func sortPurchaseInfoPriceOffers(offers []PriceOffer) {
	sort.SliceStable(offers, func(i, j int) bool {
		return comparePriceOffers(offers[i], offers[j]) < 0
	})
}
