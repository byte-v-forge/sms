package app

import (
	"context"
	"errors"

	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func runtimeCoreError(err error) *core.Error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) || status.Code(err) == codes.DeadlineExceeded {
		return core.NewError(core.CodeTimeout, "sms service request timed out", true)
	}
	if errors.Is(err, context.Canceled) || status.Code(err) == codes.Canceled {
		return core.NewError(core.CodeTimeout, "sms service request canceled", true)
	}
	if status.Code(err) == codes.Unavailable {
		return core.NewError(core.CodeSupplyUnavailable, "sms service unavailable", true)
	}
	return nil
}
