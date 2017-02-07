package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	log.SetPrefix("chmod: ")
	log.SetFlags(0)

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <mode> file ...\n", os.Args[0])
		os.Exit(1)
	}

	i, err := strconv.ParseUint(os.Args[1], 8, 32)
	if err != nil {
		log.Fatal(err)
	}
	mode := os.FileMode(i)

	for _, path := range os.Args[2:] {
		if err := os.Chmod(path, mode); err != nil {
			log.Fatal(err)
		}
	}
}
