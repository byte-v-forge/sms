package secretref

import (
	"context"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
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
