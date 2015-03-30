package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var tflag = flag.String("t", "", "type")

func main() {
	elog := log.New(os.Stderr, "umount: ", 0)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: mount [options] source target")
		flag.PrintDefaults()
	}
	var fstype string
	if *tflag != "" {
		fstype = *tflag
	} else {
		elog.Print("must set type")
		flag.Usage()
		os.Exit(1)
	}
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)
	target := flag.Arg(1)
	if err := syscall.Mount(source, target, fstype, 0, ""); err != nil {
		elog.Fatal(err)
	}
}
