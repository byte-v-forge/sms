package main

import (
	"time"

	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type runtimeServices struct {
	catalog *app.CatalogService
	order   *app.OrderService
	admin   *app.ProviderAdminService
}

func newRuntimeServices(cfg config, providers *providerspi.Registry, stores *runtimeStores, platformEventBus *natseventbus.Bus, hotStream hotstream.Bus) runtimeServices {
	clock := app.SystemClock{}
	httpTimeout := time.Duration(cfg.HTTPTimeoutSeconds) * time.Second
	orderService := app.NewOrderService(
		stores.orders,
		app.NewConfiguredProviders(providers, stores.configs, httpTimeout, cfg.ProviderHTTPProxy),
		clock,
		app.RandomIDGenerator{},
		configuredOrderEvents(platformEventBus != nil && stores.orderOutbox != nil),
		hotStream,
		stores.routeHealth,
		stores.codeSecrets,
	)
	return runtimeServices{
		catalog: app.NewCatalogService(stores.configs, providers, stores.routeHealth, httpTimeout, cfg.ProviderHTTPProxy, clock),
		order:   orderService,
		admin:   app.NewProviderAdminService(stores.configs, providers, orderService, stores.orderList, httpTimeout, cfg.ProviderHTTPProxy, hotStream),
	}
}
