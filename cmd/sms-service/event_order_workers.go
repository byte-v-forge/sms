package main

import (
	"context"
	"fmt"

	eventbusadapter "github.com/byte-v-forge/sms/internal/adapters/eventbus"
	"github.com/byte-v-forge/sms/internal/app"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
)

type orderEventWorker struct {
	name       string
	definition eventcatalog.Definition
	run        func(eventbus.Consumer) error
}

func orderEventWorkerStarts(ctx context.Context, cfg config, platformEventBus *natseventbus.Bus, orderService *app.OrderService) ([]func() error, error) {
	workers := orderEventWorkers(ctx, orderService)
	starts := make([]func() error, 0, len(workers))
	for _, worker := range workers {
		consumer, err := platformEventBus.PullWorkerForDefinition(cfg.EventStreamName, worker.definition, orderEventWorkerBatch, orderEventWorkerAckWait)
		if err != nil {
			return nil, fmt.Errorf("initialize %s worker: %w", worker.name, err)
		}
		run := worker.run
		consumerForWorker := consumer
		starts = append(starts, func() error { return run(consumerForWorker) })
	}
	return starts, nil
}

func orderEventWorkers(ctx context.Context, orderService *app.OrderService) []orderEventWorker {
	return []orderEventWorker{
		{
			name:       "SMS order acquire",
			definition: smseventcatalog.OrderAcquireRequested,
			run: func(consumer eventbus.Consumer) error {
				return eventbusadapter.RunOrderAcquireWorker(ctx, consumer, orderService)
			},
		},
		{
			name:       "SMS order poll",
			definition: smseventcatalog.OrderPollRequested,
			run: func(consumer eventbus.Consumer) error {
				return eventbusadapter.RunOrderPollWorker(ctx, consumer, orderService)
			},
		},
		{
			name:       "SMS order cancel",
			definition: smseventcatalog.OrderCancelRequested,
			run: func(consumer eventbus.Consumer) error {
				return eventbusadapter.RunOrderCancelWorker(ctx, consumer, orderService)
			},
		},
	}
}
