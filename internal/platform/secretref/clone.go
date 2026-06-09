package secretref

import (
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

func Clone(ref *commonv1.SecretRef, defaultProvider string, defaultPurpose string) *commonv1.SecretRef {
	if !Configured(ref) {
		return nil
	}
	return &commonv1.SecretRef{
		SecretId:  strings.TrimSpace(ref.GetSecretId()),
		Provider:  firstNonEmpty(ref.GetProvider(), defaultProvider),
		Purpose:   firstNonEmpty(ref.GetPurpose(), defaultPurpose),
		ExpiresAt: ref.GetExpiresAt(),
	}
}
