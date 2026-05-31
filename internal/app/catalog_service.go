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
	routeHealth      RouteHealthStore
	timeout          time.Duration
	defaultHTTPProxy string
	clock            core.Clock
}

func NewCatalogService(configs ProviderConfigStore, routeHealth RouteHealthStore, timeout time.Duration, defaultHTTPProxy string, clock core.Clock) *CatalogService {
	if clock == nil {
		clock = SystemClock{}
	}
	if routeHealth == nil {
		routeHealth = noopRouteHealthStore{}
	}
	return &CatalogService{
		configs:          configs,
		routeHealth:      routeHealth,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
		clock:            clock,
	}
}

func (s *CatalogService) ListProviders(context.Context) ([]*smsinternalv1.SmsProviderPluginDescriptor, error) {
	return listSMSProviderPluginDescriptors(), nil
}
