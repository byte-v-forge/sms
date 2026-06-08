package httpsse

import "time"

const (
	DefaultEventName   = "hotstream"
	DefaultControlName = "hotstream.control"
	DefaultHeartbeat   = 15 * time.Second
)

type ServeOptions struct {
	EventName        string
	ControlEventName string
	Heartbeat        time.Duration
	Logf             func(string, ...any)
}
