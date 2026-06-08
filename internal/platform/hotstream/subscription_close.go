package hotstream

func (s *Subscription) Close() {
	if s == nil || s.inner == nil {
		return
	}
	if s.hub != nil {
		s.hub.unsubscribe(s.inner, nil)
		return
	}
	s.inner.close(nil)
}

func (s *subscription) close(err error) {
	s.once.Do(func() {
		s.err = err
		close(s.done)
		close(s.events)
	})
}
