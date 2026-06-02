package herosms

import (
	"net/http"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

func Plugin() providerspi.Plugin {
	return providerspi.NewDefinition(providerspi.Definition{
		ProviderKey:   ProviderKey,
		DisplayName:   "HeroSMS",
		Capabilities:  providerspi.BaseCapabilities(true),
		DefaultPolicy: core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second, CancelAllowedAfter: 2 * time.Minute},
		Factory: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return New(Config{APIKey: config.GetCredentialSecret()}, client)
		},
	})
}
