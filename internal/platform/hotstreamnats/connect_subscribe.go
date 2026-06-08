package hotstreamnats

import "fmt"

func (b *Bus) subscribeConn() error {
	sub, err := b.conn.Subscribe(b.subject, b.receive)
	if err != nil {
		return fmt.Errorf("subscribe hotstream nats subject %s: %w", b.subject, err)
	}
	b.sub = sub
	if err := b.conn.Flush(); err != nil {
		b.Close()
		return fmt.Errorf("flush hotstream nats subscription: %w", err)
	}
	return nil
}
