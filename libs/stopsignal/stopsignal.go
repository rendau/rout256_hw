package stopsignal

import (
	"context"
	"os"
	"os/signal"
)

func StopSignal() <-chan struct{} {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx.Done()
}
