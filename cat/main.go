package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	elog := log.New(os.Stderr, "cat: ", 0)
	flag.Parse()
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	if len(os.Args) == 1 {
		_, err := io.Copy(stdout, bufio.NewReader(os.Stdin))
		if err != nil {
			elog.Print(err)
		}
		return
	}
	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			elog.Print(err)
		}
		fbuf := bufio.NewReader(f)
		_, err = io.Copy(stdout, fbuf)
		if err != nil {
			elog.Print(err)
		}
		f.Close()
	}
}
