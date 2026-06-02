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
