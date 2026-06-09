package app

import (
	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
)

func validateSMSCodeSecretRef(ref *commonv1.SecretRef) error {
	if err := secretref.Validate(ref); err != nil {
		return core.NewError(core.CodeValidationFailed, "sms code secret ref is invalid", false)
	}
	if ref.GetProvider() != "sms" || ref.GetPurpose() != "sms_otp" {
		return core.NewError(core.CodeValidationFailed, "sms secret ref scope mismatch", false)
	}
	return nil
}
