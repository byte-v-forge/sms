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
