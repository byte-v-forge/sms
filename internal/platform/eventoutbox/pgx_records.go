package eventoutbox

import "time"

func resolveUnixTime(now int64) int64 {
	if now <= 0 {
		return time.Now().Unix()
	}
	return now
}

func resolveBatchSize(batch int) int {
	if batch <= 0 {
		return DefaultBatch
	}
	return batch
}
