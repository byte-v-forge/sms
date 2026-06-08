package secretref

import (
	"errors"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

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
