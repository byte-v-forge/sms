package eventbus

func (expected ExpectedMessage) ValidateReceived(received ReceivedMessage) error {
	if expected.IsZero() {
		return nil
	}
	if received.Envelope == nil {
		return ErrEmptyEnvelope
	}
	if err := expected.validateSubject(received); err != nil {
		return err
	}
	if err := expected.validateMetadata(received); err != nil {
		return err
	}
	if err := expected.validatePayloadType(received.Envelope.GetPayloadType()); err != nil {
		return err
	}
	return nil
}
