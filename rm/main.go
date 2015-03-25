package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var rflag = flag.Bool("r", false, "recursive")

func main() {
	logger := log.New(os.Stderr, "rm: ", 0)
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file ...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	for _, path := range flag.Args() {
		if *rflag {
			err := os.RemoveAll(path)
			if err != nil {
				logger.Print(err)
			}
		} else {
			err := os.Remove(path)
			if err != nil {
				logger.Print(err)
			}
		}
	}
}
