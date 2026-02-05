package main

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	var i int
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			println("Hello, World!", i)
			i++
		}
	}
}
