package eventcatalog

import commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"

func Catalog() *commonv1.EventCatalog {
	definitions := All()
	out := make([]*commonv1.EventDefinition, 0, len(definitions))
	for _, definition := range definitions {
		out = append(out, definition.Proto())
	}
	return &commonv1.EventCatalog{
		StreamName:    StreamName,
		StreamSubject: StreamSubject,
		Definitions:   out,
	}
}
