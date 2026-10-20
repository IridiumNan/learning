package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("the process id: %d", syscall.Getpid())

	s := <-c
	fmt.Println("Got Signal:", s)
}
