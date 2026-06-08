package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *ProviderAdminService) GetProviderBalance(ctx context.Context, providerKey string) (core.Money, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerKey)
	if err != nil {
		return core.Money{}, err
	}
	if !config.GetEnabled() {
		return core.Money{}, core.NewError(core.CodeValidationFailed, "sms provider config is disabled", false)
	}
	provider, err := providerFromConfig(s.providers, config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}
