package grpcadapter

import (
	"context"
	"sort"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogServer) ListSmsApplications(ctx context.Context, request *smsv1.ListSmsApplicationsRequest) (*smsv1.ListSmsApplicationsResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKeys: singleProviderKey(request.GetProviderKey())})
	if err != nil {
		return &smsv1.ListSmsApplicationsResponse{Error: toProtoError(err)}, nil
	}
	items := map[string]*smsv1.SmsApplicationInfo{}
	for _, offer := range offers {
		if offer.ApplicationKey == "" {
			continue
		}
		items[offer.ApplicationKey] = &smsv1.SmsApplicationInfo{ApplicationKey: offer.ApplicationKey, DisplayName: offer.ApplicationName}
	}
	out := make([]*smsv1.SmsApplicationInfo, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetApplicationKey() < out[j].GetApplicationKey() })
	return &smsv1.ListSmsApplicationsResponse{Applications: out}, nil
}

func (s *CatalogServer) ListSmsCountries(ctx context.Context, request *smsv1.ListSmsCountriesRequest) (*smsv1.ListSmsCountriesResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKeys: singleProviderKey(request.GetProviderKey())})
	if err != nil {
		return &smsv1.ListSmsCountriesResponse{Error: toProtoError(err)}, nil
	}
	items := map[string]*smsv1.SmsCountry{}
	for _, offer := range offers {
		key := offer.CountryISO2
		if key == "" {
			key = offer.CountryCallingCode
		}
		if key == "" {
			continue
		}
		items[key] = &smsv1.SmsCountry{CountryIso2: offer.CountryISO2, Name: offer.CountryName, CountryCallingCode: offer.CountryCallingCode}
	}
	out := make([]*smsv1.SmsCountry, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetCountryIso2()+out[i].GetCountryCallingCode() < out[j].GetCountryIso2()+out[j].GetCountryCallingCode()
	})
	return &smsv1.ListSmsCountriesResponse{Countries: out}, nil
}

func singleProviderKey(providerKey string) []string {
	if providerKey == "" {
		return nil
	}
	return []string{providerKey}
}
