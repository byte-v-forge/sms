package app

import "time"

type CancelRetryError struct {
	RetryAt time.Time
}

func (e *CancelRetryError) Error() string {
	if e == nil || e.RetryAt.IsZero() {
		return "sms order cancel retry is scheduled"
	}
	return "sms order cancel retry is scheduled for " + e.RetryAt.UTC().Format(time.RFC3339)
}
