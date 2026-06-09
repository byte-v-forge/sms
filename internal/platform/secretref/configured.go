package secretref

import (
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

func Configured(ref *commonv1.SecretRef) bool {
	return strings.TrimSpace(ref.GetSecretId()) != ""
}
