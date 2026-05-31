package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"google.golang.org/protobuf/proto"
)

func cloneProviderConfig(input *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if input == nil {
		return &smsinternalv1.SmsProviderConfig{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderConfig)
}

func defaultProviderCapabilities(providerKey string) *smsinternalv1.SmsProviderCapabilities {
	if plugin, ok := smsProviderPluginByKey(providerKey); ok {
		return plugin.Capabilities()
	}
	return &smsinternalv1.SmsProviderCapabilities{}
}

func supportedProviderKey(providerKey string) bool {
	_, ok := smsProviderPluginByKey(providerKey)
	return ok
}
