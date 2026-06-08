package hotstream

func (h *Hub) unsubscribe(sub *subscription, err error) {
	if h == nil || sub == nil {
		return
	}
	h.mu.Lock()
	delete(h.subs, sub)
	h.mu.Unlock()
	sub.close(err)
}
