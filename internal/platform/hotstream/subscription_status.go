package hotstream

func (s *Subscription) Err() error {
	if s == nil || s.inner == nil {
		return nil
	}
	return s.inner.err
}

func (s *Subscription) Done() <-chan struct{} {
	if s == nil || s.inner == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}
	return s.inner.done
}
