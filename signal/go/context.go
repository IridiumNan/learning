package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var neverRead = make(chan struct{})

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Signal(syscall.SIGINT); err != nil {
		log.Fatal(err)
	}

	select {
	case <-neverRead:
		fmt.Println("ready")
	case <-ctx.Done():
		// it will be [context.Canceled] -> "context canceled"
		fmt.Println(ctx.Err())
		stop()
	}
}
