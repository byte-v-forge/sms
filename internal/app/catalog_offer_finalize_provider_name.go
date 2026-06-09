package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func finalizeOfferProviderName(offer core.RouteOffer) core.RouteOffer {
	if offer.UpstreamProviderID == "" {
		offer.UpstreamProviderID = strings.TrimSpace(offer.Route.UpstreamProviderID)
	}
	if offer.UpstreamProviderName == "" {
		offer.UpstreamProviderName = offer.UpstreamProviderID
	}
	return offer
}
