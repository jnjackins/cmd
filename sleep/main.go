package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) >= 2 {
		seconds, _ := strconv.Atoi(os.Args[1])
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}
