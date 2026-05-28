package app

import (
	"context"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

type routeOfferProvider interface {
	ListRouteOffers(context.Context, core.RouteOfferQuery) ([]core.RouteOffer, error)
}

type CatalogService struct {
	configs          ProviderConfigStore
	timeout          time.Duration
	defaultHTTPProxy string
	clock            core.Clock
}

func NewCatalogService(configs ProviderConfigStore, timeout time.Duration, defaultHTTPProxy string, clock core.Clock) *CatalogService {
	if clock == nil {
		clock = SystemClock{}
	}
	return &CatalogService{configs: configs, timeout: timeout, defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy), clock: clock}
}

func (s *CatalogService) ListProviders(context.Context) ([]*smsinternalv1.SmsProviderPluginDescriptor, error) {
	return listSMSProviderPluginDescriptors(), nil
}

func (s *CatalogService) ListPriceOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	query = normalizeOfferQuery(query)
	configs, err := s.configs.ListProviderConfigs(ctx, false, query.ProviderKey)
	if err != nil {
		return nil, err
	}
	var out []core.RouteOffer
	var lastErr error
	for _, config := range configs {
		if !config.GetEnabled() {
			continue
		}
		provider, err := providerFromConfig(config, s.timeout, s.defaultHTTPProxy)
		if err != nil {
			lastErr = err
			continue
		}
		offerProvider, ok := provider.(routeOfferProvider)
		if !ok {
			continue
		}
		offers, err := offerProvider.ListRouteOffers(ctx, query)
		if err != nil {
			lastErr = err
			continue
		}
		for _, offer := range offers {
			offer = s.finalizeOffer(config, offer)
			if !routeOfferMatches(offer, query) {
				continue
			}
			out = append(out, offer)
		}
	}
	if len(out) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return out, nil
}

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

func normalizeOfferQuery(query core.RouteOfferQuery) core.RouteOfferQuery {
	query.ApplicationKey = strings.TrimSpace(query.ApplicationKey)
	query.CountryISO2 = strings.ToUpper(strings.TrimSpace(query.CountryISO2))
	query.CountryCallingCode = strings.TrimPrefix(strings.TrimSpace(query.CountryCallingCode), "+")
	query.ProviderKey = normalizeProviderKey(query.ProviderKey)
	return query
}

func routeOfferMatches(offer core.RouteOffer, query core.RouteOfferQuery) bool {
	if query.ProviderKey != "" && !strings.EqualFold(offer.ProviderKey, query.ProviderKey) {
		return false
	}
	if query.ApplicationKey != "" && offer.ApplicationKey != "" && !strings.EqualFold(offer.ApplicationKey, query.ApplicationKey) {
		return false
	}
	if query.CountryISO2 != "" && offer.CountryISO2 != "" && !strings.EqualFold(offer.CountryISO2, query.CountryISO2) {
		return false
	}
	if query.CountryCallingCode != "" && offer.CountryCallingCode != "" && strings.TrimPrefix(offer.CountryCallingCode, "+") != query.CountryCallingCode {
		return false
	}
	return true
}
