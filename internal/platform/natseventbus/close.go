package natseventbus

func (b *Bus) Close() {
	if b == nil || b.conn == nil {
		return
	}
	b.conn.Drain()
	b.conn.Close()
}
