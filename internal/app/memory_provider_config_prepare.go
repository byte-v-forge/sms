package app

import smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"

func (s *MemoryProviderConfigStore) prepareForSave(input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := normalizeProviderConfigInput(input)
	if err != nil {
		return nil, err
	}
	if err := validateProviderConfigSupported(s.providers, config.GetProviderKey()); err != nil {
		return nil, err
	}
	return config, nil
}
