package eventcatalog

type Kind string

const (
	KindFact    Kind = "fact"
	KindCommand Kind = "command"
)

type Definition struct {
	Subject          string
	EventName        string
	EventVersion     string
	Kind             Kind
	PayloadType      string
	OwnerService     string
	ConsumerDurable  string
	Retryable        bool
	MaxDeliveries    int
	RetryDelaySecond int
}
