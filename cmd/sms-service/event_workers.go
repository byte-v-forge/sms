package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	eventbusadapter "github.com/byte-v-forge/sms/internal/adapters/eventbus"
	"github.com/byte-v-forge/sms/internal/app"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
	"golang.org/x/sync/errgroup"
)

const (
	orderEventWorkerBatch   = 10
	orderEventWorkerAckWait = 60 * time.Second
)

func configuredOrderEvents(enabled bool) app.OrderEventSink {
	if !enabled {
		return nil
	}
	return eventbusadapter.NewOrderEventRecorder("sms-service")
}

func startEventWorkers(group *errgroup.Group, ctx context.Context, cfg config, platformEventBus *natseventbus.Bus, orderOutbox *app.PostgresOrderStore, orderService *app.OrderService) error {
	if platformEventBus == nil || orderOutbox == nil {
		if strings.TrimSpace(cfg.NATSURL) != "" && orderOutbox == nil {
			log.Print("SMS platform event bus configured without PostgreSQL outbox; using in-process order flow")
		}
		return nil
	}
	starts, err := orderEventWorkerStarts(ctx, cfg, platformEventBus, orderService)
	if err != nil {
		return err
	}
	group.Go(func() error {
		return eventoutbox.RunWorker(ctx, eventoutbox.WorkerConfig{
			Name:      "sms platform event outbox",
			Processor: app.NewOrderOutboxProcessor(orderOutbox, platformEventBus),
			Logf:      log.Printf,
		})
	})
	for _, start := range starts {
		group.Go(start)
	}
	return nil
}

func orderEventWorkerStarts(ctx context.Context, cfg config, platformEventBus *natseventbus.Bus, orderService *app.OrderService) ([]func() error, error) {
	workers := []orderEventWorker{
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

type orderEventWorker struct {
	name       string
	definition eventcatalog.Definition
	run        func(eventbus.Consumer) error
}
