package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogServer) ListSmsApplications(ctx context.Context, request *smsv1.ListSmsApplicationsRequest) (*smsv1.ListSmsApplicationsResponse, error) {
	applications, err := s.service.ListApplications(ctx, core.CatalogApplicationQuery{ProviderKeys: request.GetProviderKeys()})
	if err != nil {
		return &smsv1.ListSmsApplicationsResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ListSmsApplicationsResponse{Applications: toProtoCatalogApplications(applications)}, nil
}

func (s *CatalogServer) ListSmsCountries(ctx context.Context, request *smsv1.ListSmsCountriesRequest) (*smsv1.ListSmsCountriesResponse, error) {
	countries, err := s.service.ListCountries(ctx, core.CatalogCountryQuery{ProviderKeys: request.GetProviderKeys(), ApplicationKey: request.GetApplicationKey()})
	if err != nil {
		return &smsv1.ListSmsCountriesResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ListSmsCountriesResponse{Countries: toProtoCatalogCountries(countries)}, nil
}
