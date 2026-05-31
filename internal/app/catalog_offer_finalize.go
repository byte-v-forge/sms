package app

import (
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) finalizeOffer(config *smsinternalv1.SmsProviderConfig, offer core.RouteOffer) core.RouteOffer {
	offer.ProviderKey = normalizeProviderKey(firstNonEmpty(offer.ProviderKey, offer.Route.ProviderKey, config.GetProviderKey()))
	displayName := offer.ProviderKey
	if plugin, ok := smsProviderPluginByKey(offer.ProviderKey); ok {
		displayName = plugin.DisplayName()
	}
	offer.ProviderDisplayName = firstNonEmpty(offer.ProviderDisplayName, displayName, offer.ProviderKey)
	offer.Route.ProviderKey = offer.ProviderKey
	if offer.Route.ApplicationKey == "" {
		offer.Route.ApplicationKey = offer.ApplicationKey
	}
	if offer.Route.CountryISO2 == "" {
		offer.Route.CountryISO2 = offer.CountryISO2
	}
	if offer.Route.CountryCallingCode == "" {
		offer.Route.CountryCallingCode = offer.CountryCallingCode
	}
	if offer.ApplicationKey == "" {
		offer.ApplicationKey = offer.Route.ApplicationKey
	}
	if offer.CountryISO2 == "" {
		offer.CountryISO2 = offer.Route.CountryISO2
	}
	if offer.CountryCallingCode == "" {
		offer.CountryCallingCode = offer.Route.CountryCallingCode
	}
	if offer.UpstreamProviderID == "" {
		offer.UpstreamProviderID = strings.TrimSpace(offer.Route.UpstreamProviderID)
	}
	if offer.UpstreamProviderName == "" {
		offer.UpstreamProviderName = offer.UpstreamProviderID
	}
	if offer.ObservedAt.IsZero() {
		offer.ObservedAt = s.clock.Now()
	}
	if capabilities := defaultProviderCapabilities(config.GetProviderKey()); capabilities != nil {
		offer.SupportsCancel = true
		offer.SupportsAdditionalCode = capabilities.GetSupportsAdditionalCode()
		offer.RequiresMarkMessageSent = capabilities.GetRequiresMarkMessageSent()
	}
	return offer
}
