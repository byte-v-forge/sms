package fivesim

import (
	"net/http"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

func Plugin() providerspi.Plugin {
	return providerspi.NewDefinition(providerspi.Definition{
		ProviderKey: ProviderKey,
		DisplayName: "5sim",
		Capabilities: &smsinternalv1.SmsProviderCapabilities{
			SupportsBalance:         true,
			RequiresMarkMessageSent: true,
			SupportsAdditionalCode:  true,
			SupportsCatalog:         true,
			SupportsPriceLookup:     true,
		},
		DefaultPolicy: core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second},
		Factory: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return New(Config{Token: config.GetCredentialSecret()}, client)
		},
	})
}
