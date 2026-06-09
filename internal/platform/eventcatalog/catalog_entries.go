package eventcatalog

func All() []Definition {
	return []Definition{
		SMSOrderAcquired,
		SMSCodeReceived,
		SMSOrderStatusChanged,
		DeadLetter,
	}
}

func Subjects() []string {
	return []string{StreamSubject}
}
