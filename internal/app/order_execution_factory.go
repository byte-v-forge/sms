package app

type asyncOrderEventSink interface {
	AsyncRequests() bool
}

func orderExecutionForEvents(events OrderEventSink, acquire acquireRequestRunner, cancel cancelRequestRunner) OrderExecution {
	async, ok := events.(asyncOrderEventSink)
	if ok && async.AsyncRequests() {
		return asyncOrderExecution{}
	}
	return inProcessOrderExecution{acquire: acquire, cancel: cancel}
}
