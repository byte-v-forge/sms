package main

import (
	"context"
	"log"
	"strings"
	"time"

	eventbusadapter "github.com/byte-v-forge/sms/internal/adapters/eventbus"
	"github.com/byte-v-forge/sms/internal/app"
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
