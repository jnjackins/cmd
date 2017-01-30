package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	log.SetPrefix("sleep: ")
	log.SetFlags(0)
	if len(os.Args) > 1 {
		var dur time.Duration
		seconds, err := strconv.Atoi(os.Args[1])
		if err != nil {
			if dur, err = time.ParseDuration(os.Args[1]); err != nil {
				log.Fatal(err)
			}
		} else {
			dur = time.Duration(seconds) * time.Second
		}
		time.Sleep(dur)
	}
}
