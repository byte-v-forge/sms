package app

import (
	"strings"
	"time"

	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type ProviderAdminService struct {
	configs          ProviderConfigStore
	providers        *providerspi.Registry
	orders           *OrderService
	orderDB          OrderListStore
	timeout          time.Duration
	defaultHTTPProxy string
	hot              HotStreamPublisher
}

func NewProviderAdminService(configs ProviderConfigStore, providers *providerspi.Registry, orders *OrderService, orderDB OrderListStore, timeout time.Duration, defaultHTTPProxy string, hot HotStreamPublisher) *ProviderAdminService {
	return &ProviderAdminService{
		configs:          configs,
		providers:        providers,
		orders:           orders,
		orderDB:          orderDB,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
		hot:              hot,
	}
}
