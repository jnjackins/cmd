package main

import (
	"bufio"
	"flag"
	"io"
	"os"

	"sigint.ca/die"
)

func main() {
	flag.Parse()
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	if len(os.Args) == 1 {
		io.Copy(stdout, bufio.NewReader(os.Stdin))
		return
	}
	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		die.On(err, "cat: error opening file")
		fbuf := bufio.NewReader(f)
		io.Copy(stdout, fbuf)
		f.Close()
	}
}
