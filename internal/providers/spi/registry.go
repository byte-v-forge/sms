package spi

import (
	"net/http"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

type Factory func(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)

type Plugin interface {
	Key() string
	DisplayName() string
	Capabilities() *smsinternalv1.SmsProviderCapabilities
	DefaultPolicy() core.ProviderPolicy
	NewProvider(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

type Registry struct {
	plugins map[string]Plugin
	keys    []string
}

func NormalizeKey(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
