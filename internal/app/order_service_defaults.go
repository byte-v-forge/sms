package app

import "github.com/byte-v-forge/sms/internal/core"

func orderServiceDefaults(clock core.Clock, ids core.IDGenerator, events OrderEventSink, routeHealth RouteHealthStore) (core.Clock, core.IDGenerator, OrderEventSink, RouteHealthStore) {
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
	return clock, ids, events, routeHealth
}
