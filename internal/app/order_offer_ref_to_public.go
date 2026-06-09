package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func PublicOfferRefFromRoute(route core.Route) *smsv1.SmsOfferRef {
	ref := &smsv1.SmsOfferRef{
		OfferId:     publicOfferID(route),
		ProviderKey: normalizeProviderKey(route.ProviderKey),
		Target:      publicOfferTarget(route),
		RouteRef:    publicOfferRouteRef(route),
	}
	if publicOfferRefIsZero(ref) {
		return nil
	}
	return ref
}

func publicOfferRefIsZero(ref *smsv1.SmsOfferRef) bool {
	return ref.GetProviderKey() == "" && ref.GetOfferId() == "" && targetIsZero(ref.GetTarget()) && offerRouteRefIsZero(ref.GetRouteRef())
}
