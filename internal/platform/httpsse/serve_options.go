package httpsse

import "strings"

func normalizeServeOptions(opts ServeOptions) ServeOptions {
	opts.EventName = nonEmpty(opts.EventName, DefaultEventName)
	opts.ControlEventName = nonEmpty(opts.ControlEventName, DefaultControlName)
	if opts.Heartbeat <= 0 {
		opts.Heartbeat = DefaultHeartbeat
	}
	return opts
}

func nonEmpty(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
