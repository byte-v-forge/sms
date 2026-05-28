package app

import "time"

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
