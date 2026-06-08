package eventoutbox

import "time"

func workerDelay(cfg WorkerConfig, published int) time.Duration {
	if published > 0 {
		return cfg.ActiveInterval
	}
	return cfg.Interval
}
