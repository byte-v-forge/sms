package grpcadapter

import (
	"sort"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func countryFacets(offers []core.RouteOffer) []*smsv1.SmsCountry {
	items := map[string]*smsv1.SmsCountry{}
	for _, offer := range offers {
		key := countryFacetKey(offer)
		if key == "" {
			continue
		}
		items[key] = &smsv1.SmsCountry{CountryIso2: offer.CountryISO2, Name: offer.CountryName, CountryCallingCode: offer.CountryCallingCode}
	}
	out := countryFacetValues(items)
	sort.Slice(out, func(i, j int) bool { return countrySortKey(out[i]) < countrySortKey(out[j]) })
	return out
}

func countryFacetKey(offer core.RouteOffer) string {
	if offer.CountryISO2 != "" {
		return offer.CountryISO2
	}
	return offer.CountryCallingCode
}

func countryFacetValues(items map[string]*smsv1.SmsCountry) []*smsv1.SmsCountry {
	out := make([]*smsv1.SmsCountry, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out
}

func countrySortKey(country *smsv1.SmsCountry) string {
	return country.GetCountryIso2() + country.GetCountryCallingCode()
}
