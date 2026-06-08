package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminService) ListProviderPlugins(context.Context) ([]*smsinternalv1.SmsProviderPluginDescriptor, error) {
	return s.providers.Descriptors(), nil
}
