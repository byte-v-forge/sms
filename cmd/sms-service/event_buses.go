package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/byte-v-forge/common-lib/hotstream"
	"github.com/byte-v-forge/common-lib/hotstreamnats"
	"github.com/byte-v-forge/common-lib/natseventbus"
)

func newPlatformEventBus(ctx context.Context, cfg config) (*natseventbus.Bus, func(), error) {
	if strings.TrimSpace(cfg.PlatformNATSURL) == "" {
		return nil, nil, fmt.Errorf("PLATFORM_NATS_URL is required for SMS order polling")
	}
	bus, err := natseventbus.Connect(natseventbus.Config{
		URL:        cfg.PlatformNATSURL,
		ClientName: "sms-service",
	})
	if err != nil {
		return nil, nil, err
	}
	return bus, bus.Close, nil
}

func newSMSHotStreamBus(ctx context.Context, cfg config) (hotstream.Bus, func(), error) {
	if strings.TrimSpace(cfg.PlatformNATSURL) == "" {
		return nil, nil, fmt.Errorf("PLATFORM_NATS_URL is required for SMS hotstream")
	}
	bus, err := hotstreamnats.Connect(ctx, hotstreamnats.Config{
		URL:        cfg.PlatformNATSURL,
		ClientName: "sms-service",
		Subject:    hotstream.ServiceStateSubject("sms"),
	})
	if err != nil {
		return nil, nil, err
	}
	return bus, bus.Close, nil
}
