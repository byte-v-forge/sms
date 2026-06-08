package app

import "github.com/byte-v-forge/sms/internal/core"

func NewOrderService(
	store OrderStore,
	providers []core.Provider,
	clock core.Clock,
	ids core.IDGenerator,
	events OrderEventSink,
	hot HotStreamPublisher,
	routeHealth RouteHealthStore,
	codeSecrets *SMSCodeSecretStore,
) *OrderService {
	clock, ids, events, routeHealth = orderServiceDefaults(clock, ids, events, routeHealth)
	service := &OrderService{
		store:       store,
		providers:   orderProviderMap(providers),
		routeHealth: routeHealth,
		codeSecrets: codeSecrets,
		clock:       clock,
		ids:         ids,
		events:      events,
		hot:         hot,
	}
	service.execution = orderExecutionForEvents(events, service.RunAcquireRequest, service.cancelLoadedOrder)
	return service
}
