package app

import (
	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

func cloneOrder(order core.Order) core.Order {
	order.LastError = cloneCoreError(order.LastError)
	return order
}

func cloneCoreError(err *core.Error) *core.Error {
	if err == nil {
		return nil
	}
	out := *err
	return &out
}

func cloneSMSCode(code core.SMSCode) core.SMSCode {
	code.SecretRef = cloneSecretRef(code.SecretRef)
	return code
}

func cloneSecretRef(ref *commonv1.SecretRef) *commonv1.SecretRef {
	if ref == nil {
		return nil
	}
	return proto.Clone(ref).(*commonv1.SecretRef)
}
