package eventbus

import "log"

func logger(logf LogFunc) LogFunc {
	if logf != nil {
		return logf
	}
	return log.Printf
}
