package app

import (
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type providerConfigRecord struct {
	providerKey      string
	enabled          bool
	credentialSecret string
	createdAt        time.Time
	updatedAt        time.Time
}

func (r providerConfigRecord) toProto() *smsinternalv1.SmsProviderConfig {
	return &smsinternalv1.SmsProviderConfig{
		ProviderKey:         normalizeProviderKey(r.providerKey),
		Enabled:             r.enabled,
		CredentialSecret:    r.credentialSecret,
		CredentialSecretSet: strings.TrimSpace(r.credentialSecret) != "",
		CreatedAt:           timestamppb.New(r.createdAt),
		UpdatedAt:           timestamppb.New(r.updatedAt),
	}
}
