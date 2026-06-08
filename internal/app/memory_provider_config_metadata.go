package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoryProviderConfigStore) applySaveMetadata(config *smsinternalv1.SmsProviderConfig, now *timestamppb.Timestamp) {
	existing := s.configs[config.GetProviderKey()]
	if existing != nil {
		config.CreatedAt = cloneTimestamp(existing.GetCreatedAt())
		if config.GetCredentialSecret() == "" {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
		return
	}
	config.CreatedAt = cloneTimestamp(now)
}
