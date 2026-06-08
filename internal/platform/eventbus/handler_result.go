package eventbus

import "time"

type MessageAction string

const (
	MessageActionAck  MessageAction = "ack"
	MessageActionNak  MessageAction = "nak"
	MessageActionTerm MessageAction = "term"
)

type HandlerResult struct {
	Action MessageAction
	Delay  time.Duration
	Label  string
}

func AckResult(label string) HandlerResult {
	return HandlerResult{Action: MessageActionAck, Label: label}
}

func NakResult(delay time.Duration, label string) HandlerResult {
	return HandlerResult{Action: MessageActionNak, Delay: delay, Label: label}
}

func TermResult(label string) HandlerResult {
	return HandlerResult{Action: MessageActionTerm, Label: label}
}
