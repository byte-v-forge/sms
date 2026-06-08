package handlerapi

import "github.com/byte-v-forge/sms/internal/core"

type GetNumberV2Config struct {
	ProviderName       string
	CountryLabel       string
	ProviderIDParam    string
	ProviderIDLabel    string
	ProviderIDRequired bool
	MaxPriceParam      string
}

type StatusParser func(string) (core.ProviderCodeResult, error)

type StatusActionMapper func(core.ProviderAction) (status string, expected string, err error)
