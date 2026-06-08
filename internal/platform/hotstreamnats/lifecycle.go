package hotstreamnats

func (b *Bus) Close() {
	if b == nil {
		return
	}
	if b.sub != nil {
		_ = b.sub.Unsubscribe()
	}
	if b.conn != nil {
		b.conn.Drain()
		b.conn.Close()
	}
}
