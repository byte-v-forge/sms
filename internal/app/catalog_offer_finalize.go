package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) finalizeOffer(config *smsinternalv1.SmsProviderConfig, offer core.RouteOffer) core.RouteOffer {
	offer = s.finalizeOfferProvider(config, offer)
	offer = finalizeOfferRouteDefaults(offer)
	offer = finalizeOfferProviderName(offer)
	offer = s.finalizeOfferObservedAt(offer)
	return s.finalizeOfferCapabilities(config, offer)
}
