package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupSignals() {
	c := make(chan os.Signal)
	go func() {
		for {
			_ = <-c // swallow interrupt signals
		}
	}()
	signal.Notify(c, syscall.SIGINT)
}
