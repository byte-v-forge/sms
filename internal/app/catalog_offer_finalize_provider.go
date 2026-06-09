package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) finalizeOfferProvider(config *smsinternalv1.SmsProviderConfig, offer core.RouteOffer) core.RouteOffer {
	offer.ProviderKey = normalizeProviderKey(firstNonEmpty(offer.ProviderKey, offer.Route.ProviderKey, config.GetProviderKey()))
	displayName := s.providers.DisplayName(offer.ProviderKey)
	offer.ProviderDisplayName = firstNonEmpty(offer.ProviderDisplayName, displayName, offer.ProviderKey)
	offer.Route.ProviderKey = offer.ProviderKey
	return offer
}
