package app

import (
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"google.golang.org/protobuf/proto"
)

func cloneProviderConfig(input *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if input == nil {
		return &smsinternalv1.SmsProviderConfig{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderConfig)
}

func normalizeProviderConfigInput(input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config := cloneProviderConfig(input)
	config.ProviderKey = normalizeProviderKey(config.GetProviderKey())
	if config.GetProviderKey() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	config.CredentialSecret = strings.TrimSpace(config.GetCredentialSecret())
	return config, nil
}

func validateProviderConfigSupported(providers *providerspi.Registry, providerKey string) error {
	if providers != nil && !providers.Supports(providerKey) {
		return core.NewError(core.CodeUnsupportedOperation, "unsupported sms provider", false)
	}
	return nil
}

func validateProviderConfigCredential(config *smsinternalv1.SmsProviderConfig) error {
	if config.GetEnabled() && config.GetCredentialSecret() == "" {
		return core.NewError(core.CodeValidationFailed, "credential_secret is required for enabled sms provider", false)
	}
	return nil
}

func markProviderConfigCredentialState(config *smsinternalv1.SmsProviderConfig) {
	config.CredentialSecretSet = config.GetCredentialSecret() != ""
}
