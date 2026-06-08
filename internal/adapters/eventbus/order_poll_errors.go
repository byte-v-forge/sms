package eventbusadapter

import "github.com/byte-v-forge/sms/internal/core"

func pollErrorRetryable(err error) bool {
	return coreErrorRetryableUnless(
		err,
		core.CodeOrderNotFound,
		core.CodeOrderAlreadyFinalized,
		core.CodeOrderExpired,
	)
}
