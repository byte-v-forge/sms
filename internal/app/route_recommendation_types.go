package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

const (
	defaultRouteRecommendationLimit = 20
	maxRouteRecommendationLimit     = 100
)

type RouteRecommendationQuery struct {
	Target       core.Target
	Policy       *smsv1.SmsRoutePolicy
	ProviderKeys []string
}

type RouteRecommendation struct {
	Offer core.RouteOffer
	Score int32
}
