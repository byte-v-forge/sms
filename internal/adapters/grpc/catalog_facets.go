package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogServer) ListSmsApplications(ctx context.Context, request *smsv1.ListSmsApplicationsRequest) (*smsv1.ListSmsApplicationsResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKeys: singleProviderKey(request.GetProviderKey())})
	if err != nil {
		return &smsv1.ListSmsApplicationsResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ListSmsApplicationsResponse{Applications: applicationFacets(offers)}, nil
}

func (s *CatalogServer) ListSmsCountries(ctx context.Context, request *smsv1.ListSmsCountriesRequest) (*smsv1.ListSmsCountriesResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKeys: singleProviderKey(request.GetProviderKey())})
	if err != nil {
		return &smsv1.ListSmsCountriesResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ListSmsCountriesResponse{Countries: countryFacets(offers)}, nil
}

func singleProviderKey(providerKey string) []string {
	if providerKey == "" {
		return nil
	}
	return []string{providerKey}
}
