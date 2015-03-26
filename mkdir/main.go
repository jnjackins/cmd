package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var usage = "Usage: mkdir [options] directory ..."

var pflag = flag.Bool("p", false, "create parent directories")

func main() {
	elog := log.New(os.Stderr, "mkdir: ", 0)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, usage)
		flag.PrintDefaults()
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
