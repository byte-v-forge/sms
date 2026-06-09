package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

func toProtoRouteRecommendations(recommendations []app.RouteRecommendation) []*smsv1.SmsRouteRecommendation {
	out := make([]*smsv1.SmsRouteRecommendation, 0, len(recommendations))
	for _, recommendation := range recommendations {
		out = append(out, &smsv1.SmsRouteRecommendation{
			Offer: toProtoPriceOffer(recommendation.Offer),
			Score: recommendation.Score,
		})
	}
	return out
}
