package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

func toProtoProviderLookupErrors(errors []app.ProviderLookupError) []*smsv1.SmsProviderLookupError {
	out := make([]*smsv1.SmsProviderLookupError, 0, len(errors))
	for _, providerErr := range errors {
		out = append(out, &smsv1.SmsProviderLookupError{
			ProviderKey:         providerErr.ProviderKey,
			ProviderDisplayName: providerErr.ProviderDisplayName,
			Error:               toProtoError(providerErr.Err),
		})
	}
	return out
}
