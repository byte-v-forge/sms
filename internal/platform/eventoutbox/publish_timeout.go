package eventoutbox

import "time"

const defaultPublishTimeout = 10 * time.Second

func publishTimeout(options PublishOptions) time.Duration {
	if options.PublishTimeout > 0 {
		return options.PublishTimeout
	}
	return defaultPublishTimeout
}
