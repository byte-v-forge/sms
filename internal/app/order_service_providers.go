package app

import "github.com/byte-v-forge/sms/internal/core"

func orderProviderMap(providers []core.Provider) map[string]core.Provider {
	index := make(map[string]core.Provider, len(providers))
	for _, provider := range providers {
		index[provider.Key()] = provider
	}
	return index
}
