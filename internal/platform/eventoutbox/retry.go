package eventoutbox

import "time"

func DefaultRetryDelay(attempt int32) time.Duration {
	switch {
	case attempt <= 1:
		return 5 * time.Second
	case attempt == 2:
		return 15 * time.Second
	case attempt == 3:
		return 30 * time.Second
	case attempt <= 6:
		return time.Minute
	default:
		return 5 * time.Minute
	}
}

func retryDelay(options PublishOptions, attempt int32) time.Duration {
	if options.RetryDelay != nil {
		return options.RetryDelay(attempt)
	}
	return DefaultRetryDelay(attempt)
}
