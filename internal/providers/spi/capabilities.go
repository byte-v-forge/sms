package spi

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"google.golang.org/protobuf/proto"
)

func BaseCapabilities(catalog bool) *smsinternalv1.SmsProviderCapabilities {
	return &smsinternalv1.SmsProviderCapabilities{
		SupportsBalance:         true,
		RequiresMarkMessageSent: true,
		SupportsAdditionalCode:  true,
		SupportsCatalog:         catalog,
		SupportsPriceLookup:     catalog,
	}
}

func cloneCapabilities(input *smsinternalv1.SmsProviderCapabilities) *smsinternalv1.SmsProviderCapabilities {
	if input == nil {
		return &smsinternalv1.SmsProviderCapabilities{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderCapabilities)
}
