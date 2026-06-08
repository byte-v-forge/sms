package hotstreamnats

import "context"

func closeBusOnContext(ctx context.Context, bus *Bus) {
	go func() {
		<-ctx.Done()
		bus.Close()
	}()
}
