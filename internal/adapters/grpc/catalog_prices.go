package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogServer) ListSmsPriceOffers(ctx context.Context, request *smsv1.ListSmsPriceOffersRequest) (*smsv1.ListSmsPriceOffersResponse, error) {
	result, err := s.service.ListPriceOffersDetailed(ctx, core.RouteOfferQuery{
		ApplicationKey:     request.GetApplicationKey(),
		CountryISO2:        request.GetCountryIso2(),
		CountryCallingCode: request.GetCountryCallingCode(),
		ProviderKeys:       request.GetProviderKeys(),
	})
	response := &smsv1.ListSmsPriceOffersResponse{
		Offers:         toProtoPriceOffers(result.Offers),
		ProviderErrors: toProtoProviderLookupErrors(result.ProviderErrors),
	}
	if err != nil {
		response.Error = toProtoError(err)
	}
	return response, nil
}

func (s *CatalogServer) RecommendSmsRoutes(ctx context.Context, request *smsv1.RecommendSmsRoutesRequest) (*smsv1.RecommendSmsRoutesResponse, error) {
	target := request.GetTarget()
	recommendations, err := s.service.RecommendRoutes(ctx, app.RouteRecommendationQuery{
		Target: core.Target{
			ApplicationKey:     target.GetApplicationKey(),
			CountryISO2:        target.GetCountryIso2(),
			CountryCallingCode: target.GetCountryCallingCode(),
		},
		Policy:       request.GetPolicy(),
		ProviderKeys: request.GetProviderKeys(),
	})
	if err != nil {
		return &smsv1.RecommendSmsRoutesResponse{Error: toProtoError(err)}, nil
	}
	out := make([]*smsv1.SmsRouteRecommendation, 0, len(recommendations))
	for _, recommendation := range recommendations {
		out = append(out, &smsv1.SmsRouteRecommendation{
			Offer: toProtoPriceOffer(recommendation.Offer),
			Score: recommendation.Score,
		})
	}
	return &smsv1.RecommendSmsRoutesResponse{Recommendations: out}, nil
}

func toProtoProviderLookupErrors(errors []app.ProviderLookupError) []*smsv1.SmsProviderLookupError {
	out := make([]*smsv1.SmsProviderLookupError, 0, len(errors))
	for _, providerErr := range errors {
		out = append(out, &smsv1.SmsProviderLookupError{
			ProviderKey:         providerErr.ProviderKey,
			ProviderDisplayName: providerErr.ProviderDisplayName,
			Error:               toProtoError(providerErr.Err),
		})
	}
	return out
}
