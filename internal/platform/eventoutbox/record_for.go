package eventoutbox

import (
	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"google.golang.org/protobuf/proto"
)

func NewRecordFor(
	definition eventcatalog.Definition,
	event proto.Message,
	metadata *commonv1.EventMetadata,
	attributes map[string]string,
) (Record, error) {
	message, err := definition.NewMessage(event, metadata, attributes)
	if err != nil {
		return Record{}, err
	}
	return NewRecord(message)
}
