package app

import (
	"context"

	"github.com/byte-v-forge/common-lib/eventoutbox"
	"github.com/byte-v-forge/sms/internal/core"
)

type OrderService struct {
	store       OrderStore
	hot         HotStreamPublisher
	providers   map[string]core.Provider
	routeHealth RouteHealthStore
	codeSecrets *SMSCodeSecretStore
	clock       core.Clock
	ids         core.IDGenerator
	events      OrderEventSink
}

type OrderStore interface {
	Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error
	Get(ctx context.Context, orderID string) (core.Order, error)
	Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error
	RecordCode(ctx context.Context, order core.Order, code core.SMSCode, events ...eventoutbox.Record) error
}

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
	index := make(map[string]core.Provider, len(providers))
	for _, provider := range providers {
		index[provider.Key()] = provider
	}
	if clock == nil {
		clock = SystemClock{}
	}
	if ids == nil {
		ids = RandomIDGenerator{}
	}
	if events == nil {
		events = noopOrderEventSink{}
	}
	if routeHealth == nil {
		routeHealth = noopRouteHealthStore{}
	}
	return &OrderService{
		store:       store,
		providers:   index,
		routeHealth: routeHealth,
		codeSecrets: codeSecrets,
		clock:       clock,
		ids:         ids,
		events:      events,
		hot:         hot,
	}
}
