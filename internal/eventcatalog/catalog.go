package eventcatalog

import commoncatalog "github.com/byte-v-forge/common-lib/eventcatalog"

var (
	OrderAcquireRequested = commoncatalog.SMSOrderAcquireRequested
	OrderPollRequested    = commoncatalog.SMSOrderPollRequested
	OrderCancelRequested  = commoncatalog.SMSOrderCancelRequested
)
