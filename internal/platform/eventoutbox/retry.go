package eventoutbox

import (
	"strings"
	"time"
)

const defaultPublishTimeout = 10 * time.Second

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

func TruncateError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.TrimSpace(err.Error())
	if len(message) <= 1000 {
		return message
	}
	return message[:1000]
}

func publishTimeout(options PublishOptions) time.Duration {
	if options.PublishTimeout > 0 {
		return options.PublishTimeout
	}
	return defaultPublishTimeout
}

func retryDelay(options PublishOptions, attempt int32) time.Duration {
	if options.RetryDelay != nil {
		return options.RetryDelay(attempt)
	}
	return DefaultRetryDelay(attempt)
}

func optionNow(options PublishOptions) time.Time {
	if options.Now != nil {
		return options.Now()
	}
	return time.Now()
}
