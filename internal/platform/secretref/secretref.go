package secretref

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/hashx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WriteRequest struct {
	SecretID  string
	Provider  string
	Purpose   string
	Value     string
	ExpiresAt time.Time
}

type Writer interface {
	WriteSecret(ctx context.Context, req WriteRequest) (*commonv1.SecretRef, error)
}

type Resolver interface {
	ResolveSecret(ctx context.Context, ref *commonv1.SecretRef) (string, error)
}

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

func Validate(ref *commonv1.SecretRef) error {
	if !Configured(ref) {
		return errors.New("secret_id is required")
	}
	if strings.TrimSpace(ref.GetProvider()) == "" {
		return errors.New("secret provider is required")
	}
	if strings.TrimSpace(ref.GetPurpose()) == "" {
		return errors.New("secret purpose is required")
	}
	return nil
}

func StableID(prefix string, parts ...string) string {
	prefix = cleanSegment(prefix)
	if prefix == "" {
		prefix = "secret"
	}
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			clean = append(clean, value)
		}
	}
	if len(clean) == 0 {
		return prefix
	}
	return fmt.Sprintf("%s-%s", prefix, hashx.ShortSHA256(hashx.StableParts(clean...), 24))
}

func Display(ref *commonv1.SecretRef) string {
	if !Configured(ref) {
		return ""
	}
	provider := strings.TrimSpace(ref.GetProvider())
	purpose := strings.TrimSpace(ref.GetPurpose())
	id := strings.TrimSpace(ref.GetSecretId())
	if len(id) > 12 {
		id = id[:12]
	}
	switch {
	case provider != "" && purpose != "":
		return provider + "/" + purpose + "/" + id
	case provider != "":
		return provider + "/" + id
	default:
		return id
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func cleanSegment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(" ", "-", "_", "-", "/", "-", ":", "-").Replace(value)
	value = strings.Trim(value, "-")
	return value
}
