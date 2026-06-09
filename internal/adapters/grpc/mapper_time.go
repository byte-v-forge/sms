package grpcadapter

import (
	"time"

	"github.com/byte-v-forge/sms/internal/app"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoTime(value time.Time) *timestamppb.Timestamp { return app.PublicTime(value) }

func protoDuration(value *durationpb.Duration) time.Duration {
	if value == nil {
		return 0
	}
	return value.AsDuration()
}
