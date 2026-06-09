package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func (s *CatalogServer) ListSmsPriceOffers(ctx context.Context, request *smsv1.ListSmsPriceOffersRequest) (*smsv1.ListSmsPriceOffersResponse, error) {
	result, err := s.service.ListPriceOffersDetailed(ctx, priceOfferQueryFromRequest(request))
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
	recommendations, err := s.service.RecommendRoutes(ctx, recommendationQueryFromRequest(request))
	if err != nil {
		return &smsv1.RecommendSmsRoutesResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.RecommendSmsRoutesResponse{Recommendations: toProtoRouteRecommendations(recommendations)}, nil
}
