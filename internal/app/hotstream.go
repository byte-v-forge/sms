package app

import "github.com/byte-v-forge/sms/internal/platform/hotstream"

const (
	SMSHotStreamSource        = "sms-service"
	SMSOrderResource          = "sms.order"
	SMSProviderConfigResource = "sms.provider_config"
	SMSOrderUpdatedEvent      = "sms.order.updated"
	SMSProviderConfigUpdated  = "sms.provider_config.updated"
	SMSProviderConfigDeleted  = "sms.provider_config.deleted"
)

type HotStreamPublisher = hotstream.Publisher
