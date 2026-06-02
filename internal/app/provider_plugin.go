package app

import (
	"fmt"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

func NewConfiguredProviders(providers *providerspi.Registry, configs ProviderConfigStore, timeout time.Duration, defaultHTTPProxy string) []core.Provider {
	descriptors := providers.Descriptors()
	out := make([]core.Provider, 0, len(descriptors))
	for _, descriptor := range descriptors {
		out = append(out, NewConfiguredProvider(providers, descriptor.GetProviderKey(), configs, timeout, defaultHTTPProxy))
	}
	return out
}

func unsupportedSMSProvider(providerKey string) error {
	return core.NewError(core.CodeUnsupportedOperation, fmt.Sprintf("unsupported sms provider %q", providerKey), false)
}
