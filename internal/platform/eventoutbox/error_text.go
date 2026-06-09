package eventoutbox

import "strings"

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
