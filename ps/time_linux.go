package main

import "syscall"

func ticksPerSecond() int {
	var t syscall.Timex
	syscall.Adjtimex(&t)
	// 1e6 microseconds per second
	return int(1e6 / t.Tick)
}
