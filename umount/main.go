package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	elog := log.New(os.Stderr, "umount: ", 0)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: umount mountpoint ...")
	}
	for _, path := range flag.Args() {
		if err := unix.Unmount(path, 0); err != nil {
			elog.Print(err)
		}
	}
}
