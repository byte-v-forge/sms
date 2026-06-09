package eventbusadapter

import (
	"fmt"
	"time"
)

func eventTimeSuffix(value time.Time) string {
	return fmt.Sprintf("%d", value.UnixNano())
}
