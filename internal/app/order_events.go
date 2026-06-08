package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

type OrderEventSink interface {
	OrderAcquireRequested(context.Context, core.Order, core.Route, string) (eventoutbox.Record, error)
	OrderAcquired(context.Context, core.Order) (eventoutbox.Record, error)
	OrderPollRequested(context.Context, core.Order, string) (eventoutbox.Record, error)
	OrderCancelRequested(context.Context, core.Order, string, string) (eventoutbox.Record, error)
	CodeReceived(context.Context, core.Order, core.SMSCode) (eventoutbox.Record, error)
	OrderStatusChanged(context.Context, core.Order, core.OrderStatus) (eventoutbox.Record, error)
}

type noopOrderEventSink struct{}

func (noopOrderEventSink) OrderAcquireRequested(context.Context, core.Order, core.Route, string) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}

func (noopOrderEventSink) OrderAcquired(context.Context, core.Order) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}

func (noopOrderEventSink) OrderPollRequested(context.Context, core.Order, string) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}

func (noopOrderEventSink) OrderCancelRequested(context.Context, core.Order, string, string) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}

func (noopOrderEventSink) CodeReceived(context.Context, core.Order, core.SMSCode) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}

func (noopOrderEventSink) OrderStatusChanged(context.Context, core.Order, core.OrderStatus) (eventoutbox.Record, error) {
	return eventoutbox.Record{}, nil
}
