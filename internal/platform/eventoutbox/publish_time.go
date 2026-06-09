package eventoutbox

import "time"

func optionNow(options PublishOptions) time.Time {
	if options.Now != nil {
		return options.Now()
	}
	return time.Now()
}
