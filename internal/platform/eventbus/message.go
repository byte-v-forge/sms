package eventbus

import (
	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/proto"
)

const ProtobufContentType = "application/x-protobuf"

type Message struct {
	Subject    string
	Event      proto.Message
	Metadata   *commonv1.EventMetadata
	Extensions map[string]string
}
