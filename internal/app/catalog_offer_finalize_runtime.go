package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) finalizeOfferObservedAt(offer core.RouteOffer) core.RouteOffer {
	if offer.ObservedAt.IsZero() {
		offer.ObservedAt = s.clock.Now()
	}
	return offer
}

func (s *CatalogService) finalizeOfferCapabilities(config *smsinternalv1.SmsProviderConfig, offer core.RouteOffer) core.RouteOffer {
	if capabilities := s.providers.Capabilities(config.GetProviderKey()); capabilities != nil {
		offer.SupportsCancel = true
		offer.SupportsAdditionalCode = capabilities.GetSupportsAdditionalCode()
		offer.RequiresMarkMessageSent = capabilities.GetRequiresMarkMessageSent()
	}
	return offer
}
