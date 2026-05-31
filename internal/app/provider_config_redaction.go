package app

import (
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"google.golang.org/protobuf/proto"
)

func RedactProviderConfig(config *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if config == nil {
		return nil
	}
	redacted := proto.Clone(config).(*smsinternalv1.SmsProviderConfig)
	redacted.CredentialSecretSet = strings.TrimSpace(config.GetCredentialSecret()) != "" || config.GetCredentialSecretSet()
	redacted.CredentialSecret = ""
	return redacted
}
