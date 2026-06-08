package app

import (
	"fmt"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
)

func smsCodeSecretRef(orderID string, receivedAt time.Time, expiresAt time.Time) *commonv1.SecretRef {
	secretID := secretref.StableID("sms-code", orderID, fmt.Sprintf("%d", receivedAt.UnixNano()))
	return secretref.New("sms", "sms_otp", secretID, expiresAt)
}
