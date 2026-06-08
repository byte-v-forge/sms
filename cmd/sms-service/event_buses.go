package main

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/byte-v-forge/sms/internal/platform/hotstreamnats"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
)

func newPlatformEventBus(_ context.Context, cfg config) (*natseventbus.Bus, func(), error) {
	if strings.TrimSpace(cfg.NATSURL) == "" {
		return nil, noopClose, nil
	}
	bus, err := natseventbus.Connect(natseventbus.Config{
		URL:        cfg.NATSURL,
		ClientName: "sms-service",
	})
	if err != nil {
		return nil, nil, err
	}
	return bus, bus.Close, nil
}

func newSMSHotStreamBus(ctx context.Context, cfg config) (hotstream.Bus, func(), error) {
	if strings.TrimSpace(cfg.NATSURL) == "" {
		return hotstream.NewHub(hotstream.DefaultBufferSize), noopClose, nil
	}
	bus, err := hotstreamnats.ConnectService(ctx, hotstreamnats.ServiceConfig{
		URL:        cfg.NATSURL,
		ClientName: "sms-service",
		Service:    "sms",
	})
	if err != nil {
		return nil, nil, err
	}
	return bus, bus.Close, nil
}

func noopClose() {}
