package eventbus

import commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"

func NewEnvelope(message Message) (*commonv1.EventEnvelope, error) {
	subject, err := envelopeSubject(message.Subject)
	if err != nil {
		return nil, err
	}
	payload, payloadType, err := envelopePayload(message)
	if err != nil {
		return nil, err
	}
	if err := ValidateMetadata(message.Metadata); err != nil {
		return nil, err
	}
	return &commonv1.EventEnvelope{
		Metadata:        message.Metadata,
		Subject:         subject,
		PayloadType:     payloadType,
		Payload:         payload,
		DataContentType: ProtobufContentType,
		Extensions:      cloneExtensions(message.Extensions),
	}, nil
}
