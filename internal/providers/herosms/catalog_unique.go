package herosms

import "strings"

func uniqueHeroSMSOffers(offers []PriceOffer) []PriceOffer {
	seen := map[string]struct{}{}
	out := make([]PriceOffer, 0, len(offers))
	for _, offer := range offers {
		key := offer.CountryID + "\x00" + offer.UpstreamServiceKey + "\x00" + offer.Operator + "\x00" + offer.Price.AmountDecimal
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, offer)
	}
	return out
}

func uniqueHeroSMSStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if _, exists := seen[value]; value == "" || exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
