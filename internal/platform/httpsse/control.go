package httpsse

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

func control(kind observabilityv1.HotStreamControlKind, message string) *observabilityv1.HotStreamControlEvent {
	return &observabilityv1.HotStreamControlEvent{Kind: kind, Message: message, OccurredAt: timestamppb.Now()}
}

func protoJSON(message proto.Message) string {
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(message)
	if err != nil {
		fallback, _ := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_ERROR, err.Error()))
		return string(fallback)
	}
	return string(data)
}
