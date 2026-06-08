package spi

import (
	"net/http"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

type Definition struct {
	ProviderKey   string
	DisplayName   string
	Capabilities  *smsinternalv1.SmsProviderCapabilities
	DefaultPolicy core.ProviderPolicy
	Factory       Factory
}

type definitionPlugin struct{ definition Definition }

func NewDefinition(definition Definition) Plugin {
	definition.ProviderKey = NormalizeKey(definition.ProviderKey)
	definition.DisplayName = strings.TrimSpace(definition.DisplayName)
	if definition.Capabilities == nil {
		definition.Capabilities = &smsinternalv1.SmsProviderCapabilities{}
	}
	return definitionPlugin{definition: definition}
}

func (p definitionPlugin) Key() string { return p.definition.ProviderKey }

func (p definitionPlugin) DisplayName() string { return p.definition.DisplayName }

func (p definitionPlugin) Capabilities() *smsinternalv1.SmsProviderCapabilities {
	return cloneCapabilities(p.definition.Capabilities)
}

func (p definitionPlugin) DefaultPolicy() core.ProviderPolicy { return p.definition.DefaultPolicy }

func (p definitionPlugin) NewProvider(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
	if p.definition.Factory == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms provider factory is not configured", false)
	}
	return p.definition.Factory(config, client)
}
