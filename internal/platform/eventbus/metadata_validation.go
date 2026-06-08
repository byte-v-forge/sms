package eventbus

import (
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
)

func ValidateMetadata(metadata *commonv1.EventMetadata) error {
	if metadata == nil {
		return ErrEmptyEventMetadata
	}
	if strings.TrimSpace(metadata.GetId()) == "" {
		return ErrEmptyEventID
	}
	if strings.TrimSpace(metadata.GetType()) == "" {
		return ErrEmptyEventType
	}
	if strings.TrimSpace(metadata.GetVersion()) == "" {
		return ErrEmptyEventVersion
	}
	if metadata.GetTime() == nil || !metadata.GetTime().IsValid() {
		return ErrEmptyEventTime
	}
	if strings.TrimSpace(metadata.GetSource()) == "" {
		return ErrEmptySource
	}
	if strings.TrimSpace(metadata.GetIdempotencyKey()) == "" {
		return ErrEmptyIdempotencyKey
	}
	return nil
}
