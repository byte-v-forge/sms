package eventbus

func EventID(message ReceivedMessage) string {
	if message.Envelope == nil || message.Envelope.GetMetadata() == nil {
		return ""
	}
	return message.Envelope.GetMetadata().GetId()
}
