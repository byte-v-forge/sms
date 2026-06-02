package main

import (
	"github.com/byte-v-forge/sms/internal/providers/fivesim"
	"github.com/byte-v-forge/sms/internal/providers/herosms"
	"github.com/byte-v-forge/sms/internal/providers/smsbower"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

func newProviderRegistry() (*providerspi.Registry, error) {
	return providerspi.NewRegistry(
		fivesim.Plugin(),
		herosms.Plugin(),
		smsbower.Plugin(),
	)
}
