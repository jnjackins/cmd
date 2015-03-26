package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var pflag = flag.Bool("p", false, "create parent directories")

func main() {
	elog := log.New(os.Stderr, "mkdir: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: mkdir [options] directory ...")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	mkdir := os.Mkdir
	if *pflag {
		mkdir = os.MkdirAll
	}
	for _, dir := range args {
		err := mkdir(dir, 0777)
		if err != nil {
			elog.Print(err)
		}
	}
}
