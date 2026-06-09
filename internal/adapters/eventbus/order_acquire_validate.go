package eventbusadapter

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func validateAcquireRequest(request *smsinternalv1.SmsOrderAcquireRequest) error {
	return validateOrderID(request.GetOrderId())
}

func acquireErrorRetryable(err error) bool {
	return coreErrorRetryableUnless(
		err,
		core.CodeValidationFailed,
		core.CodeUnsupportedOperation,
		core.CodeInsufficientBalance,
		core.CodeOrderExpired,
		core.CodeOrderAlreadyFinalized,
	)
}
