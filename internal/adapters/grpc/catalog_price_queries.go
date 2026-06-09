package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

func priceOfferQueryFromRequest(request *smsv1.ListSmsPriceOffersRequest) core.RouteOfferQuery {
	return core.RouteOfferQuery{
		ApplicationKey:     request.GetApplicationKey(),
		CountryISO2:        request.GetCountryIso2(),
		CountryCallingCode: request.GetCountryCallingCode(),
		ProviderKeys:       request.GetProviderKeys(),
	}
}

func recommendationQueryFromRequest(request *smsv1.RecommendSmsRoutesRequest) app.RouteRecommendationQuery {
	target := request.GetTarget()
	return app.RouteRecommendationQuery{
		Target: core.Target{
			ApplicationKey:     target.GetApplicationKey(),
			CountryISO2:        target.GetCountryIso2(),
			CountryCallingCode: target.GetCountryCallingCode(),
		},
		Policy:       request.GetPolicy(),
		ProviderKeys: request.GetProviderKeys(),
	}
}
