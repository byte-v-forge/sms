package eventoutbox

import "github.com/byte-v-forge/sms/internal/platform/eventbus"

func NewRecord(message eventbus.Message) (Record, error) {
	envelope, err := eventbus.NewEnvelope(message)
	if err != nil {
		return Record{}, err
	}
	return recordFromEnvelope(envelope)
}
