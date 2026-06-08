package secretref

import (
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(provider string, purpose string, secretID string, expiresAt time.Time) *commonv1.SecretRef {
	secretID = strings.TrimSpace(secretID)
	if secretID == "" {
		return nil
	}
	ref := &commonv1.SecretRef{
		SecretId: secretID,
		Provider: strings.TrimSpace(provider),
		Purpose:  strings.TrimSpace(purpose),
	}
	if !expiresAt.IsZero() {
		ref.ExpiresAt = timestamppb.New(expiresAt)
	}
	return ref
}

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

func Configured(ref *commonv1.SecretRef) bool {
	return strings.TrimSpace(ref.GetSecretId()) != ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
